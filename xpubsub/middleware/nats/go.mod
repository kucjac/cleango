module github.com/kucjac/cleango/xpubsub/pubsubmiddleware/natsmiddleware

go 1.16

require (
	github.com/kucjac/cleango v0.0.27
	github.com/nats-io/nats.go v1.13.0
	gocloud.dev v0.24.0
	gocloud.dev/pubsub/natspubsub v0.24.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)

replace github.com/kucjac/cleango => ../../../
