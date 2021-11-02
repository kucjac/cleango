package xservice

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/pkg/xlog"
)

// RunService runs the service and locks the thread while running given service.
func RunService(ctx context.Context, s RunnerCloser) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	errorChan := make(chan error, 1)

	// Start the service.
	go func(s Runner, cancel context.CancelFunc) {
		var err error
		if err = s.Run(); err != nil {
			if cgerrors.Code(err) != cgerrors.CodeCanceled {
				xlog.Errorf("ListenAndServe failed: %v", err)
				errorChan <- err
			} else {
				cancel()
			}
		}
	}(s, cancel)

	select {
	case <-ctx.Done():
		xlog.Infof("Service context had finished.")
	case sig := <-quit:
		xlog.Infof("Received Signal: '%s'. Shutdown Server begins...", sig.String())
	case err := <-errorChan:
		// The error from the server running.
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.Close(ctx); err != nil {
		xlog.Errorf("Server shutdown failed: %v", err)
		return err
	}
	xlog.Info("Server had shutdown successfully.")
	return nil
}
