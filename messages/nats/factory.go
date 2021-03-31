package nats

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"

	"github.com/kucjac/cleango/messages"
	"github.com/kucjac/cleango/messages/codec"
	"github.com/kucjac/cleango/messages/pubsub"
	"github.com/kucjac/cleango/xlog"
)

var (
	_ pubsub.SubscriberFactory = (*factory)(nil)
	_ pubsub.PublisherFactory  = (*factory)(nil)
)

// factory is the nats publisher and subscriber factory adapter.
type factory struct {
	Conn   stan.Conn
	Codec  codec.Codec
	Logger watermill.LoggerAdapter
}

// PublisherFactory gets the messages.PublisherFactory type from given factory.
func (f *factory) PublisherFactory() pubsub.PublisherFactory {
	return f
}

// SubscriberFactory gets the messages.SubscriberFactory type from given factory.
func (f *factory) SubscriberFactory() pubsub.SubscriberFactory {
	return f
}

// NewFactory creates new nats.factory for the messages.SubscriberFactory and messages.PublisherFactory interfaces.
func NewFactory(c codec.Codec, config nats.StanConnConfig, logger xlog.Logger) (pubsub.Factory, error) {
	s, err := nats.NewStanConnection(&config)
	if err != nil {
		return nil, err
	}
	return &factory{
		Conn:   s,
		Codec:  c,
		Logger: messages.NewLoggerAdapter(logger),
	}, nil
}

// NewSubscriber creates new subscriber
func (f *factory) NewSubscriber(options ...pubsub.SubscriptionOption) (pubsub.Subscriber, error) {
	so := &pubsub.SubscriptionOptions{}
	for _, option := range options {
		option(so)
	}
	cfg := nats.StreamingSubscriberSubscriptionConfig{
		Unmarshaler: nats.GobMarshaler{},
	}
	cfg.AckWaitTimeout = so.AckWaitTimeout
	cfg.CloseTimeout = so.CloseTimeout
	cfg.SubscribersCount = so.SubscribersCount
	cfg.DurableName = so.DurableName
	cfg.QueueGroup = so.QueueGroup

	switch so.StartAt {
	case pubsub.StartPosition_TimeDeltaStart:
		cfg.StanSubscriptionOptions = append(cfg.StanSubscriptionOptions, stan.StartAtTime(so.StartTime))
	case pubsub.StartPosition_SequenceStart:
		cfg.StanSubscriptionOptions = append(cfg.StanSubscriptionOptions, stan.StartAtSequence(so.StartSequence))
	case pubsub.StartPosition_LastReceived:
		cfg.StanSubscriptionOptions = append(cfg.StanSubscriptionOptions, stan.StartWithLastReceived())
	case pubsub.StartPosition_First:
		cfg.StanSubscriptionOptions = append(cfg.StanSubscriptionOptions, stan.DeliverAllAvailable())
	}

	if so.MaxInflight != 0 {
		cfg.StanSubscriptionOptions = append(cfg.StanSubscriptionOptions, stan.MaxInflight(so.MaxInflight))
	}

	return nats.NewStreamingSubscriberWithStanConn(f.Conn, cfg, f.Logger)
}

// NewPublisher creates new messages publisher.
// Implements messages.PublisherFactory
func (f *factory) NewPublisher() (pubsub.Publisher, error) {
	natsPub, err := nats.NewStreamingPublisherWithStanConn(f.Conn, nats.StreamingPublisherPublishConfig{Marshaler: nats.GobMarshaler{}}, f.Logger)
	if err != nil {
		return nil, err
	}
	return &publisher{pub: natsPub}, nil
}

// publisher is the nats publisher implementation.
type publisher struct {
	c   codec.Codec
	pub *nats.StreamingPublisher
}

// Publish implements messages.Publisher interface.
func (p *publisher) Publish(topic string, messages ...*message.Message) error {
	return p.pub.Publish(topic, messages...)
}

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
	return p.pub.Publish(topic, m)
}

// Close implements messages.Publisher interface.
func (p *publisher) Close() error {
	return p.pub.Close()
}
