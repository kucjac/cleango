package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/kucjac/cleango/eventstore"
)

// EventSubscription contains parameters required for handling event subscriptions.
type EventSubscription struct {
	Topic   string
	Handler eventstore.EventHandler
	Options []SubscriptionOption
}

// MessageHandler is an interface used for handling the message.
type MessageHandler interface {
	Handle(ctx context.Context, message *message.Message)
}

// MessageSubscription is message subscription option.
type MessageSubscription struct {
	Topic   string
	Handler MessageHandler
	Options []SubscriptionOption
}
