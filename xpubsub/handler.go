package xpubsub

import (
	"context"

	"gocloud.dev/pubsub"
)

// Handler is an interface used by the Mux that is responsible for handling the Message.
type Handler interface {
	Handle(ctx context.Context, m *pubsub.Message) error
}

// HandlerFunc is a function used to handle messages.
type HandlerFunc func(ctx context.Context, m *pubsub.Message) error

func (h HandlerFunc) Handle(ctx context.Context, m *pubsub.Message) error {
	return h(ctx, m)
}
