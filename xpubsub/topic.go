package xpubsub

import (
	"context"

	"gocloud.dev/pubsub"
)

//go:generate mockgen -destination=mockpubsub/topic_gen.go -package=mockpubsub . Topic

// Compile time check if pubsub.Topic implements Topic interface.
var _ Topic = (*pubsub.Topic)(nil)

// Topic is the interface implementation of the pubsub.Topic.
// It allows to easily mockup the pubsub.Topic by replacing direct implementattion
type Topic interface {
	Send(ctx context.Context, m *pubsub.Message) error
	Shutdown(ctx context.Context) error
	ErrorAs(err error, i interface{}) bool
	As(i interface{}) bool
}
