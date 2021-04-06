package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

// MessageHandler is an interface used for handling the message.
type MessageHandler interface {
	Handle(ctx context.Context, topic string, message *message.Message) error
}

// MessageSubscription is message subscription option.
type MessageSubscription struct {
	Topic   string
	Options []SubscriptionOption
}
