package codec

import (
	"bytes"
	"encoding/gob"
)

// GOB gets std gob codec.
func GOB() Codec {
	return &gobCodec{}
}

var _ Codec = (*gobCodec)(nil)

type gobCodec struct{}

// Name gets the name of the codec.
func (j *gobCodec) Name() string {
	return "gob"
}

// Marshal implements codec.Marshaler interface.
func (j *gobCodec) Marshal(input interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(input); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal implements codec.Marshaler interface
func (j *gobCodec) Unmarshal(data []byte, output interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(output)
}

// Encoding returns codec encoding.
func (j *gobCodec) Encoding() string {
	return encodingGOB
}

const encodingGOB = "application/gob"

// EncodingGOB returns json encoding.
func EncodingGOB() string {
	return encodingGOB
}
