package pubsub

// Message is the event message interface.
type Message interface {
	MessageTopic() string
}
