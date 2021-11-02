package xpubsub

import (
	"context"

	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/pkg/xlog"
	"gocloud.dev/pubsub"
	"google.golang.org/grpc/metadata"
)

//go:generate mockgen -destination=mock/topic_gen.go -package=pubsubmock . Topic

// Compile time check if pubsub.Topic implements Topic interface.
var _ Topic = (*pubsub.Topic)(nil)

// Topic is the interface implementation of the pubsub.Topic.
// It allows to easily mockup the pubsub.Topic by replacing direct implementation.
type Topic interface {
	Send(ctx context.Context, m *pubsub.Message) error
	Shutdown(ctx context.Context) error
	ErrorAs(err error, i interface{}) bool
	As(i interface{}) bool
}

// TopicPublisher is a structure responsible for publishing new topic messages.
type TopicPublisher struct {
	Topic Topic
	Codec codec.Codec
}

// Publish prepares pubsub.Message with the context stored metadata and publishes into given topic.
func (t *TopicPublisher) Publish(ctx context.Context, msg MessageTyper) error {
	// Marshal the message body.
	body, err := t.Codec.Marshal(msg)
	if err != nil {
		xlog.WithContext(ctx).Errorf("marshaling '%s' message failed: %v", msg.MessageType(), err)
		return err
	}

	// Define the message.
	m := pubsub.Message{Body: body}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		m.Metadata = map[string]string{}
		for k, v := range md {
			if len(v) != 0 {
				m.Metadata[k] = v[0]
			}
		}
	}

	// Publish the message into given topic.
	if err = t.Topic.Send(ctx, &m); err != nil {
		xlog.WithContext(ctx).Errorf("sending '%s' message failed: %v", msg.MessageType(), err)
		return err
	}
	return nil
}
