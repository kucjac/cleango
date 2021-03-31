package xservice

// RunnerCloser is the interface used as one of the service ports.
type RunnerCloser interface {
	Run() error
	Close() error
}
