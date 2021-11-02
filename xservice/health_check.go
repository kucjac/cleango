package xservice

//go:generate enumer -text -trimprefix=HealthCheckStatus -transform=snake_upper -type=HealthCheckStatus -output health_check_status_string.go

// HealthCheckStatus gets the status of the health check.
type HealthCheckStatus int

// Health check predefined statuses.
const (
	HealthCheckStatusUnknown HealthCheckStatus = iota
	HealthCheckStatusServing
	HealthCheckStatusNotServing
	HealthCheckStatusServiceUnknown
)
