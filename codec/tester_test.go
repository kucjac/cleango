package codec

import (
	"testing"

	"github.com/kucjac/cleango/codec/internal/pb"
	"google.golang.org/protobuf/proto"
)

type testMessage struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (t *testMessage) MarshalProto() ([]byte, error) {
	msg := pb.TestMessage{
		Id:   t.ID,
		Name: t.Name,
	}
	return proto.Marshal(&msg)
}

func (t *testMessage) UnmarshalProto(data []byte) error {
	var msg pb.TestMessage
	if err := proto.Unmarshal(data, &msg); err != nil {
		return err
	}
	t.ID = msg.Id
	t.Name = msg.Name
	return nil
}

func TestTester(t *testing.T) {
	TestMessagesEncoding(t, []Codec{Proto(), JSON(), GOB()}, testMessage{ID: "123", Name: "MyName"})
}
