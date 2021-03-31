package messages

import (
	"github.com/ThreeDotsLabs/watermill/message"
	uuid "github.com/satori/go.uuid"

	"github.com/kucjac/cleango/messages/codec"
	"github.com/kucjac/cleango/meta"
)

// EncodeMessage gets the event message
func EncodeMessage(c codec.Codec, e Messager, metadata meta.Meta) (*message.Message, error) {
	payload, err := c.Marshal(e)
	if err != nil {
		return nil, err
	}
	msg := message.NewMessage(uuid.NewV4().String(), payload)
	if metadata != nil {
		msg.Metadata = message.Metadata(metadata)
	}
	return msg, nil
}

// Messager is the event message interface.
type Messager interface {
	MessageTopic() string
}
