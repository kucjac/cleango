package codec

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var _ Codec = (*protoCodec)(nil)

type protoCodec struct{}

// Marshal implements codec.Marshaler interface.
func (j *protoCodec) Marshal(input interface{}) ([]byte, error) {
	msg, ok := input.(proto.Message)
	if !ok {
		return nil, errors.New("protobuf input is not a proto.Message")
	}
	return proto.Marshal(msg)
}

// Unmarshal implements codec.Marshaler interface
func (j *protoCodec) Unmarshal(data []byte, output interface{}) error {
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
