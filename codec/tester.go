package codec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessagesEncoding tests the input messages if the values encoding persists the same.
func TestMessagesEncoding(t *testing.T, codecs []Codec, messages ...interface{}) {
	for _, msg := range messages {
		t.Run(fmt.Sprintf("%T", msg), func(t *testing.T) {
			for _, e := range codecs {
				t.Run(e.Encoding(), func(t *testing.T) {
					msgPtrTp := reflect.New(reflect.TypeOf(msg))
					msgPtrTp.Elem().Set(reflect.ValueOf(msg))
					msgPtr := msgPtrTp.Interface()

					msg2PtrTp := reflect.New(reflect.TypeOf(msg))
					msg2Ptr := msg2PtrTp.Interface()

					data, err := e.Marshal(msgPtr)
					require.NoError(t, err)

					err = e.Unmarshal(data, msg2Ptr)
					require.NoError(t, err)

					msg2 := reflect.ValueOf(msg2Ptr).Elem().Interface()

					assert.Equal(t, msg, msg2)
				})
			}
		})
	}
}
