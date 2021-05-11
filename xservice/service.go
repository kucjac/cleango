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
	name    string
	runners []RunnerCloser
	closers []Closer
	err     error
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
	if err := s.run(); err != nil {
		s.err = err
	}
	return nil
}

// Start locks the thread and serves the service runners. The resultant channel closes when the service is finished.
func (s *Service) Start(ctx context.Context) <-chan struct{} {
	r := make(chan struct{}, 1)
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
				xlog.Errorf("ListenAndServe failed: %v", err)
				errorChan <- err
			} else {
				cancel()
			}
		}
	}(cancel)

	select {
	case <-ctx.Done():
		xlog.Infof("Service context had finished.")
	case sig := <-quit:
		xlog.Infof("Received Signal: '%s'. Shutdown Server begins...", sig.String())
	case err := <-errorChan:
		// The error from the server running.
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := s.Close(); err != nil {
		xlog.Errorf("Server shutdown failed: %v", err)
		return err
	}
	xlog.Info("Server had shutdown successfully.")
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
		xlog.Errorf("%s - context deadline exceeded: %v", s.name, ctx.Err())
		return ctx.Err()
	case e := <-errChan:
		xlog.Errorf("%s error: %v", s.name, e)
		return e
	case <-waitChan:
		xlog.Debugf("%s successfully started", s.name)
	}
	return nil
}

// Close closes all connection within provided context.
func (s *Service) Close() error {
	xlog.Infof("Closing %s service and its runners...", s.name)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	wg := &sync.WaitGroup{}
	waitChan := make(chan struct{})
	jobs := s.closeJobsCreator(ctx, wg)

	errChan := make(chan error)
	for job := range jobs {
		xlog.Debugf("Closing: %T", job)
		s.closeCloser(job, wg, errChan)
	}

	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-ctx.Done():
		xlog.Errorf("Close %s - context deadline exceeded: %v", s.name, ctx.Err())
		return ctx.Err()
	case e := <-errChan:
		xlog.Debugf("Close %s error: %v", e, s.name)
		return e
	case <-waitChan:
		xlog.Debugf("Closed %s all repositories with success", s.name)
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

func (s *Service) closeCloser(closer Closer, wg *sync.WaitGroup, errChan chan<- error) {
	go func() {
		defer wg.Done()
		if err := closer.Close(); err != nil {
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