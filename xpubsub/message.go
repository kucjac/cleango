package xpubsub

import (
	"github.com/kucjac/cleango/codec"
	"gocloud.dev/pubsub"
)

// MarshalMessage is a function that marshals input message structure into the pubsub.Message body.
func MarshalMessage(c codec.Codec, msg interface{}) (*pubsub.Message, error) {
	data, err := c.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &pubsub.Message{Body: data}, nil
}

// UnmarshalMessage unmarshals message body into provided destination using given codec.
func UnmarshalMessage(c codec.Codec, msg *pubsub.Message, dst interface{}) error {
	return c.Unmarshal(msg.Body, dst)
}

// MessageTyper is an interface used by the event messages that defines the type of the message.
type MessageTyper interface {
	MessageType() string
}
