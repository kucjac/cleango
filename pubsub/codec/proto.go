package codec

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// Proto gets std protobuf codec.
func Proto() Codec {
	return &protoCodec{}
}

var _ Codec = (*protoCodec)(nil)

// ProtoMarshaler is an interface used by the protobuf codec used to custom marshal proto message.
type ProtoMarshaler interface {
	MarshalProto() ([]byte, error)
}

// ProtoUnmarshaler is an interface used for custom unmarshalling proto message.
type ProtoUnmarshaler interface {
	UnmarshalProto([]byte) error
}

type protoCodec struct{}

// Marshal implements codec.Marshaler interface.
func (j *protoCodec) Marshal(input interface{}) ([]byte, error) {
	if mp, ok := input.(ProtoMarshaler); ok {
		return mp.MarshalProto()
	}
	msg, ok := input.(proto.Message)
	if !ok {
		return nil, errors.New("protobuf input is not a proto.Message")
	}
	return proto.Marshal(msg)
}

// Unmarshal implements codec.Marshaler interface
func (j *protoCodec) Unmarshal(data []byte, output interface{}) error {
	up, ok := output.(ProtoUnmarshaler)
	if ok {
		return up.UnmarshalProto(data)
	}
	msg, ok := output.(proto.Message)
	if !ok {
		return errors.New("protobuf output is not a proto.Message")
	}
	return proto.Unmarshal(data, msg)
}

// Encoding returns codec encoding.
func (j *protoCodec) Encoding() string {
	return encodingProto
}

const encodingProto = "application/proto"

// EncodingProto returns json encoding.
func EncodingProto() string {
	return encodingProto
}
