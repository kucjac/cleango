package chanbus

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	uuid "github.com/satori/go.uuid"

	"github.com/kucjac/cleango/messages"
	"github.com/kucjac/cleango/messages/codec"
	"github.com/kucjac/cleango/messages/pubsub"
	"github.com/kucjac/cleango/xlog"
)

var (
	_ pubsub.Factory           = (*factory)(nil)
	_ pubsub.PublisherFactory  = (*factory)(nil)
	_ pubsub.SubscriberFactory = (*factory)(nil)
)

// New creates a new factory.
func New(c codec.Codec, cfg gochannel.Config, logger xlog.Logger) pubsub.Factory {
	logAdapter := messages.NewLoggerAdapter(logger)
	return &factory{
		gc:         gochannel.NewGoChannel(cfg, logAdapter),
		c:          c,
		logAdapter: logAdapter,
	}
}

type factory struct {
	gc         *gochannel.GoChannel
	c          codec.Codec
	logAdapter watermill.LoggerAdapter
}

func (f *factory) PublisherFactory() pubsub.PublisherFactory {
	return f
}

func (f *factory) SubscriberFactory() pubsub.SubscriberFactory {
	return f
}

func (f *factory) NewSubscriber(_ ...pubsub.SubscriptionOption) (pubsub.Subscriber, error) {
	return f.gc, nil
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
func (p *publisher) PublishMessage(topic string, msg interface{}, options ...pubsub.PublishOption) error {
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

func (p *publisher) Close() error {
	return p.g.Close()
}

func (f *factory) NewPublisher() (pubsub.Publisher, error) {
	return &publisher{g: f.gc}, nil
}
