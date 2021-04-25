package pubsub

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/kucjac/cleango/meta"
	"github.com/kucjac/cleango/pubsub/codec"
)

// Factory is the interface used for both the publisher and subscriber factory.
type Factory interface {
	PublisherFactory() PublisherFactory
	SubscriberFactory() SubscriberFactory
}

// Publisher is the interface used by the
type Publisher interface {
	// Publish publishes provided messages to given topic.
	//
	// Publish can be synchronous or asynchronous - it depends on the implementation.
	//
	// Most publishers implementations don't support atomic publishing of messages.
	// This means that if publishing one of the messages fails, the next messages will not be published.
	//
	// Publish must be thread safe.
	Publish(topic string, messages ...*message.Message) error
	// PublishMessage encodes provided input message msg
	PublishMessage(topic string, msg interface{}, options ...PublishOption) error
	// Close should flush unsent messages, if publisher is async.
	Close() error
}

// PublishOption is an option function that changes the 'PublishMessage' message.
type PublishOption func(m *message.Message)

// PublishWithMeta publishes message with given metadata.
func PublishWithMeta(metadata meta.Meta) PublishOption {
	return func(m *message.Message) {
		m.Metadata = message.Metadata(metadata)
	}
}

// PublishWithUUID publishes message with given identifier.
func PublishWithUUID(uuid string) PublishOption {
	return func(m *message.Message) {
		m.UUID = uuid
	}
}

// Subscriber is the interface used for registering to given subscription.
type Subscriber interface {
	// Subscribe returns output channel with messages from provided topic.
	// Channel is closed, when Close() was called on the subscriber.
	//
	// To receive the next message, `Ack()` must be called on the received message.
	// If message processing failed and message should be redelivered `Nack()` should be called.
	//
	// When provided ctx is cancelled, subscriber will close subscribe and close output channel.
	// Provided ctx is set to all produced messages.
	// When Nack or Ack is called on the message, context of the message is canceled.
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	// Close closes all subscriptions with their output channels and flush offsets etc. when needed.
	Close() error
}

// SubscribeInitializer is the interface used for initializing the subscriber.
type SubscribeInitializer interface {
	// SubscribeInitialize can be called to initialize subscribe before consume.
	// When calling Subscribe before Publish, SubscribeInitialize should be not required.
	//
	// Not every Pub/Sub requires this initialize and it may be optional for performance improvements etc.
	// For detailed SubscribeInitialize functionality, please check Pub/Subs godoc.
	//
	// Implementing SubscribeInitialize is not obligatory.
	SubscribeInitialize(topic string) error
}

// PublisherFactory is the factory for the publishers.
type PublisherFactory interface {
	NewPublisher(c codec.Codec) (Publisher, error)
}

// SubscriberFactory is the factory for the subscribers.
type SubscriberFactory interface {
	NewSubscriber(options ...SubscriptionOption) (Subscriber, error)
}

type SubscriptionOptions struct {
	// QueueGroup is the NATS Streaming queue group.
	//
	// All subscriptions with the same queue name (regardless of the connection they originate from)
	// will form a queue group. Each message will be delivered to only one subscriber per queue group,
	// using queuing semantics.
	//
	// It is recommended to set it with DurableName.
	// For non durable queue subscribers, when the last member leaves the group,
	// that group is removed. A durable queue group (DurableName) allows you to have all members leave
	// but still maintain state. When a member re-joins, it starts at the last position in that group.
	//
	// When QueueGroup is empty, subscribe without QueueGroup will be used.
	QueueGroup string

	// DurableName is the NATS streaming durable name.
	//
	// Subscriptions may also specify a “durable name” which will survive client restarts.
	// Durable subscriptions cause the server to track the last acknowledged message
	// sequence number for a client and durable name. When the client restarts/resubscribes,
	// and uses the same client ID and durable name, the server will resume delivery beginning
	// with the earliest unacknowledged message for this durable subscription.
	//
	// Doing this causes the NATS Streaming server to track
	// the last acknowledged message for that ClientID + DurableName.
	DurableName string

	// SubscribersCount determines wow much concurrent subscribers should be started.
	SubscribersCount int

	// CloseTimeout determines how long subscriber will wait for Ack/Nack on close.
	// When no Ack/Nack is received after CloseTimeout, subscriber will be closed.
	CloseTimeout time.Duration

	// How long subscriber should wait for Ack/Nack. When no Ack/Nack was received, message will be redelivered.
	// It is mapped to stan.AckWait option.
	AckWaitTimeout time.Duration
	// These subscription options are
	MaxInflight   int
	StartAt       uint32
	StartSequence uint64
	StartTime     time.Time
}

// SubscriptionOption is an option used as a subscription option function. Changes the settings for given subscription.
type SubscriptionOption func(o *SubscriptionOptions)

// SubQueueGroup creates a subscribe queue for given subscriber.
func SubQueueGroup(name string) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.QueueGroup = name
	}
}

// SubDurableName creates a durable name for given message.
func SubDurableName(durableName string) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.DurableName = durableName
	}
}

// SubsCount sets the SubsCount option for subscription.
func SubsCount(count int) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.SubscribersCount = count
	}
}

// SubCloseTimeout sets the CloseTimeout option for subscription.
func SubCloseTimeout(closeTO time.Duration) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.CloseTimeout = closeTO
	}
}

// SubAckWaitTimeout sets the AckWaitTimeout option for subscription.
func SubAckWaitTimeout(awt time.Duration) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.AckWaitTimeout = awt
	}
}

// SubMaxInflight is an Option to set the maximum number of messages the cluster will send
// without an ACK.
func SubMaxInflight(m int) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.MaxInflight = m
	}
}

// SubscribeStartAt sets the desired start position for the message stream.
func SubscribeStartAt(sp uint32) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartAt = sp
	}
}

const (
	StartPosition_NewOnly        uint32 = 0
	StartPosition_LastReceived   uint32 = 1
	StartPosition_TimeDeltaStart uint32 = 2
	StartPosition_SequenceStart  uint32 = 3
	StartPosition_First          uint32 = 4
)

// SubStartAtSequence sets the desired start sequence position and state.
func SubStartAtSequence(seq uint64) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartAt = StartPosition_SequenceStart
		o.StartSequence = seq
	}
}

// SubStartAtTime sets the desired start time position and state.
func SubStartAtTime(start time.Time) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartTime = start
		o.StartAt = StartPosition_TimeDeltaStart
	}
}

// SubStartAtTimeDelta sets the desired start time position and state using the delta.
func SubStartAtTimeDelta(ago time.Duration) SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartTime = time.Now().Add(-ago)
		o.StartAt = StartPosition_TimeDeltaStart
	}
}

// SubStartWithLastReceived is a helper function to set start position to last received.
func SubStartWithLastReceived() SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartAt = StartPosition_LastReceived
	}
}

// SubDeliverAllAvailable will deliver all messages available.
func SubDeliverAllAvailable() SubscriptionOption {
	return func(o *SubscriptionOptions) {
		o.StartAt = StartPosition_First
	}
}
