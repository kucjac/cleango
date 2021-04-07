package xservice

// RunnerCloser is the interface used as one of the service ports.
type RunnerCloser interface {
	Run() error
	Close() error
}

// Transactioner is an interface used as a transaction base, responsible for committing and rolling back the transaction.
type Transactioner interface {
	Commit() error
	Rollback() error
}

// Runner is an interface used for the services which allows to run it.
type Runner interface {
	Run() error
}
