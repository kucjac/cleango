module github.com/kucjac/cleango/xpubsub/pubsubmiddleware/natsmiddleware

go 1.16

require (
	github.com/kucjac/cleango v0.0.24
	github.com/nats-io/nats.go v1.12.1
	gocloud.dev v0.24.0
	gocloud.dev/pubsub/natspubsub v0.24.0
)

replace github.com/kucjac/cleango => ../../../
