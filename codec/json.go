package codec

import (
	"encoding/json"
)

// JSON gets std json codec implementation.
func JSON() Codec {
	return &jsonCodec{}
}

var _ Codec = (*jsonCodec)(nil)

type jsonCodec struct{}

// Name implements Codec interface.
func (j *jsonCodec) Name() string {
	return "json"
}

// Marshal implements codec.Marshaler interface.
func (j *jsonCodec) Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

// Unmarshal implements codec.Marshaler interface
func (j *jsonCodec) Unmarshal(data []byte, output interface{}) error {
	return json.Unmarshal(data, output)
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
