package xservice

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"
)

// New creates a new service implementation that allows to add sub runners (ports).
func New(name string) *Service {
	if name == "" {
		name = "Service"
	}
	return &Service{name: name}
}

// Service is an implementation of the service that allows to start and close ports and other RunnerCloser.
type Service struct {
	name              string
	runners           []RunnerCloser
	closers           []Closer
	err               error
	RunWithoutRunners bool
}

// With sets the sub runner (port) that would be started and closed along with the service.
func (s *Service) With(sub RunnerCloser) {
	s.runners = append(s.runners, sub)
}

// WithCloser adds the closer to the given service, which would be closed along with the service.
func (s *Service) WithCloser(closer Closer) {
	s.closers = append(s.closers, closer)
}

// Run establish connection for all dialers in the service.
func (s *Service) Run() error {
	if err := s.canRun(); err != nil {
		s.err = err
		return s.err
	}
	if err := s.run(); err != nil {
		s.err = err
	}
	return s.err
}

// Start locks the thread and serves the service runners. The resultant channel closes when the service is finished.
func (s *Service) Start(ctx context.Context) <-chan struct{} {
	r := make(chan struct{}, 1)
	if err := s.canRun(); err != nil {
		s.err = err
		close(r)
		return r
	}

	go func(r chan struct{}) {
		if err := s.serve(ctx); err != nil {
			s.err = err
		}
		close(r)
	}(r)
	return r
}

// Error gets the service resultant error.
func (s *Service) Error() error {
	return s.err
}


// Close closes all connection within provided context.
func (s *Service) Close(ctx context.Context) error {
	xlog.Infof("Closing Service '%s' and its runners...", s.name)

	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc()

	wg := &sync.WaitGroup{}
	waitChan := make(chan struct{})
	jobs := s.closeJobsCreator(ctx, wg)

	errChan := make(chan error)
	for job := range jobs {
		xlog.Debugf("Closing: %T", job)
		s.closeCloser(ctx, job, wg, errChan)
	}

	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-ctx.Done():
		xlog.Errorf("Close Service '%s' - context deadline exceeded: %v", s.name, ctx.Err())
		return ctx.Err()
	case e := <-errChan:
		xlog.Errorf("Close Service '%s' error: %v", e, s.name)
		return e
	case <-waitChan:
		xlog.Infof("Closed Service '%s' all repositories with success", s.name)
	}
	return nil
}

func (s *Service) serve(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	errorChan := make(chan error, 1)

	// Run the service.
	go func(cancel context.CancelFunc) {
		var err error
		if err = s.run(); err != nil {
			if cgerrors.Code(err) != cgerrors.ErrorCode_Canceled {
				xlog.Errorf("running service: %s failed: %v", s.name, err)
				errorChan <- err
			} else {
				cancel()
			}
		}
	}(cancel)

	select {
	case <-ctx.Done():
		xlog.Infof("Service: '%s' context had finished.", s.name)
	case sig := <-quit:
		xlog.Infof("Received Signal: '%s'. Shutdown Service: '%s' begins...", sig.String(), s.name)
	case err := <-errorChan:
		// The error from the server running.
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := s.Close(ctx); err != nil {
		xlog.Errorf("Service: '%s' shutdown failed: %v", s.name, err)
		return err
	}
	xlog.Infof("Service: '%s' had shutdown successfully.", s.name)
	return nil
}


func (s *Service) run() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	wg := &sync.WaitGroup{}
	waitChan := make(chan struct{})

	jobs := s.runJobsCreator(ctx, wg)
	// Create error channel.
	errChan := make(chan error)

	// Runner to all repositories.
	for job := range jobs {
		s.runSub(job, wg, errChan)
	}
	// Create wait group channel finish function.
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-ctx.Done():
		xlog.Errorf("Service '%s' - context deadline exceeded: %v", s.name, ctx.Err())
		return ctx.Err()
	case e := <-errChan:
		xlog.Errorf("Service '%s' error: %v", s.name, e)
		return e
	case <-waitChan:
		xlog.Infof("Service '%s' successfully started", s.name)
	}
	return nil
}



func (s *Service) closeJobsCreator(ctx context.Context, wg *sync.WaitGroup) <-chan Closer {
	out := make(chan Closer)
	go func() {
		defer close(out)
		for _, port := range s.runners {
			wg.Add(1)
			select {
			case out <- port:
			case <-ctx.Done():
				return
			}
		}
		for _, closer := range s.closers {
			wg.Add(1)
			select {
			case out <- closer:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func (s *Service) closeCloser(ctx context.Context, closer Closer, wg *sync.WaitGroup, errChan chan<- error) {
	go func() {
		defer wg.Done()
		if err := closer.Close(ctx); err != nil {
			errChan <- err
			return
		}
	}()
}

func (s *Service) runJobsCreator(ctx context.Context, wg *sync.WaitGroup) <-chan RunnerCloser {
	out := make(chan RunnerCloser)
	go func() {
		defer close(out)
		// Iterate over file stores.
		for _, port := range s.runners {
			wg.Add(1)
			select {
			case out <- port:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func (s *Service) runSub(runner RunnerCloser, wg *sync.WaitGroup, errChan chan<- error) {
	go func() {
		if err := runner.Run(); err != nil {
			errChan <- err
			return
		}
	}()
	time.Sleep(time.Millisecond * 200)
	wg.Done()
}

func (s *Service) canRun() error {
	if len(s.runners) > 0 || s.RunWithoutRunners {
		return nil
	}
	return cgerrors.ErrInternalf("nothing to run in: %s service", s.name)
}
