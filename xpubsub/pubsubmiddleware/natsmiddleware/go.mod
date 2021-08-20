module github.com/kucjac/cleango/xpubsub/pubsubmiddleware/natsmiddleware

go 1.16

require (
	github.com/kucjac/cleango v0.0.24
	github.com/nats-io/nats-server/v2 v2.2.6 // indirect
	github.com/nats-io/nats.go v1.11.0
	gocloud.dev v0.23.0
	gocloud.dev/pubsub/natspubsub v0.23.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
)

replace github.com/kucjac/cleango => ../../../
