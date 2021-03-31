package errors

import (
	er "errors"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestFromError(t *testing.T) {
	var err error
	err = ErrNotFoundf("%s", "example")
	merr := FromError(err)
	if merr.Code != uint32(codes.NotFound) {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
	err = er.New(err.Error())
	merr = FromError(err)
	if merr.Code != uint32(codes.NotFound) {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}

}

func TestEqual(t *testing.T) {
	err1 := ErrNotFound("msg1")
	err2 := ErrNotFound("msg2")

	if !Equal(err1, err2) {
		t.Fatal("errors must be equal")
	}

	err3 := er.New("my test err")
	if Equal(err1, err3) {
		t.Fatal("errors must be not equal")
	}

}

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			Id:     "test",
			Code:   uint32(codes.Internal),
			Detail: "Internal error",
		},
	}

	for _, e := range testData {
		ne := New(e.Id, e.Detail, codes.Code(e.Code))

		if e.Error() != ne.Error() {
			t.Fatalf("Expected %s got %s", e.Error(), ne.Error())
		}

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.Id != e.Id {
			t.Fatalf("Expected %s got %s", e.Id, pe.Id)
		}

		if pe.Detail != e.Detail {
			t.Fatalf("Expected %s got %s", e.Detail, pe.Detail)
		}

		if pe.Code != e.Code {
			t.Fatalf("Expected %d got %d", e.Code, pe.Code)
		}
	}
}
