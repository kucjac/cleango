package codec

import (
	jsoniter "github.com/json-iterator/go"
)

var _ Codec = (*jsonCodec)(nil)

type jsonCodec struct{}

// Marshal implements codec.Marshaler interface.
func (j *jsonCodec) Marshal(input interface{}) ([]byte, error) {
	return jsoniter.Marshal(input)
}

// Unmarshal implements codec.Marshaler interface
func (j *jsonCodec) Unmarshal(data []byte, output interface{}) error {
	return jsoniter.Unmarshal(data, output)
}

// Encoding returns codec encoding.
func (j *jsonCodec) Encoding() string {
	return encodingJSON
}

const encodingJSON = "application/json"

// EncodingJSON returns json encoding.
func EncodingJSON() string {
	return encodingJSON
}
