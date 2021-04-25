package chanps

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	pubsub2 "github.com/kucjac/cleango/pubsub"
	uuid "github.com/satori/go.uuid"

	"github.com/kucjac/cleango/errors"
	"github.com/kucjac/cleango/pubsub/codec"
	"github.com/kucjac/cleango/xlog"
)

var (
	_ pubsub2.Factory           = (*factory)(nil)
	_ pubsub2.PublisherFactory  = (*factory)(nil)
	_ pubsub2.SubscriberFactory = (*factory)(nil)
)

// New creates a new factory.
func New(cfg gochannel.Config, logger xlog.Logger) pubsub2.Factory {
	logAdapter := pubsub2.NewLoggerAdapter(logger)
	return &factory{
		gc:         gochannel.NewGoChannel(cfg, logAdapter),
		logAdapter: logAdapter,
	}
}

type factory struct {
	gc         *gochannel.GoChannel
	logAdapter watermill.LoggerAdapter
}

// PublisherFactory implements pubsub.Factory interface.
func (f *factory) PublisherFactory() pubsub2.PublisherFactory {
	return f
}

// SubscriberFactory implements pubsub.Factory interface.
func (f *factory) SubscriberFactory() pubsub2.SubscriberFactory {
	return f
}

// NewSubscriber implements pubsub.SubscriberFactory interface.
func (f *factory) NewSubscriber(_ ...pubsub2.SubscriptionOption) (pubsub2.Subscriber, error) {
	return f.gc, nil
}

// NewPublisher creates new channel based publisher.
// Implements pubsub.PublisherFactory.
func (f *factory) NewPublisher(c codec.Codec) (pubsub2.Publisher, error) {
	if c == nil {
		return nil, errors.ErrInternal("no codec provided for chan publisher")
	}
	return &publisher{g: f.gc, c: c}, nil
}

type publisher struct {
	g *gochannel.GoChannel
	c codec.Codec
}

// Publish implements messages.Publisher interface.
func (p *publisher) Publish(topic string, messages ...*message.Message) error {
	return p.g.Publish(topic, messages...)
}

// PublishMessage encodes provided input message and publishes on provided topic.
// Implements messages.Publisher interface.
func (p *publisher) PublishMessage(topic string, msg interface{}, options ...pubsub2.PublishOption) error {
	payload, err := p.c.Marshal(msg)
	if err != nil {
		return err
	}
	m := message.NewMessage("", payload)
	for _, option := range options {
		option(m)
	}
	if m.UUID == "" {
		m.UUID = uuid.NewV4().String()
	}
	return p.g.Publish(topic, m)
}

// Close closes given publisher connection.
// Implements pubsub.Publisher interface.
func (p *publisher) Close() error {
	return p.g.Close()
}
