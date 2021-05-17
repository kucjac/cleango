package xservice

import (
	"context"
)

// RunnerCloser is the interface used as one of the service ports.
type RunnerCloser interface {
	Run() error
	Close(ctx context.Context) error
}

// Transactioner is an interface used as a transaction base, responsible for committing and rolling back the transaction.
type Transactioner interface {
	Commit() error
	Rollback() error
}

// Runner is an interface used for the services which allows to runSub it.
type Runner interface {
	Run() error
}

// Closer is an interface used for the services which allows to got closed.
type Closer interface {
	Close(ctx context.Context) error
}

// Starter is an interface used for the services that starts and blocks it's thread.
type Starter interface {
	Start(ctx context.Context) <-chan struct{}
}
