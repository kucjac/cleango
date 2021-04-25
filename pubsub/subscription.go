package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

// MessageHandler is an interface used for handling the message.
type MessageHandler interface {
	Handle(ctx context.Context, topic string, message *message.Message) error
}

type Handler interface {
	Handle(m *message.Message)
}

// HandlerFunc is a function used to handle messages.
type HandlerFunc func(m *message.Message)

func (h HandlerFunc) Handle(m *message.Message) {
	h(m)
}

// MessageSubscription is message subscription option.
type MessageSubscription struct {
	Topic       string
	HandlerFunc HandlerFunc
	Options     []SubscriptionOption
}
