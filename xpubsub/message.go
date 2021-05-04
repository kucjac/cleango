package xpubsub

import (
	"context"

	"gocloud.dev/pubsub"
)

// Message is a pubsub message wrapper that contains a context.
type Message struct {
	ctx context.Context
	*pubsub.Message
}

// WithContext sets up the context of the message.
func (m *Message) WithContext(ctx context.Context) *Message {
	if ctx == nil {
		panic("nil context")
	}
	m2 := new(Message)
	m2.ctx = ctx
	return m2
}

// Context receives the context of the message.
func (m *Message) Context() context.Context {
	if m.ctx == nil {
		return context.Background()
	}
	return m.ctx
}
