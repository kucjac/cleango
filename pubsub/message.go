package pubsub

// MessageTopicer is the interface used for the messages that has it's topic defined and returned in the MessageTopic method.
type MessageTopicer interface {
	MessageTopic() string
}
