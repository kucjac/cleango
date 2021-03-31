package codec

import (
	"fmt"
)

// Codec is an interface for all codecs to implement the
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Encoding() string
}

var codecMap = map[string]Codec{
	encodingJSON:  &jsonCodec{},
	encodingGOB:   &gobCodec{},
	encodingProto: &protoCodec{},
}

// Register registers the codec in the mapping.
func Register(codec Codec) error {
	if _, ok := codecMap[codec.Encoding()]; ok {
		return fmt.Errorf("codec: '%s' already registered", codec.Encoding())
	}
	codecMap[codec.Encoding()] = codec
	return nil
}

// Get returns a codec for provided encoding name.
func Get(encoding string) (Codec, bool) {
	c, ok := codecMap[encoding]
	return c, ok
}
