package cgerrors

import (
	"errors"
	"testing"
)

func TestFromError(t *testing.T) {
	err := ErrNotFoundf("%s", "example")
	t.Run("cgerrors.Error", func(t *testing.T) {
		err2 := FromError(err)
		if !Equal(err, err2) {
			t.Fatalf("invalid conversation %v != %v", err, err2)
		}
	})

	t.Run("Simple", func(t *testing.T) {
		stdErr := errors.New(err.Error())
		err2 := FromError(stdErr)
		if !Equal(err2, err) {
			t.Fatalf("invalid conversation %v != %v", err, err2)
		}
	})

	t.Run("GRPCStatus", func(t *testing.T) {
		// Get the GRPC Status from the message.
		grpcErr := ToGRPCError(err)

		err2 := FromError(grpcErr)

		if !Equal(err, err2) {
			t.Errorf("invalid grpc status conversion: GRPC Err: %v, Result: %v, Expected: %v", grpcErr, err2, err)
		}
	})
}

func TestEqual(t *testing.T) {
	err1 := ErrNotFound("msg1")
	err2 := ErrNotFound("msg2")

	if !Equal(err1, err2) {
		t.Fatal("errors must be equal")
	}

	err3 := errors.New("my test err")
	if Equal(err1, err3) {
		t.Fatal("errors must be not equal")
	}
}

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			ID:     "test",
			Code:   CodeInternal,
			Detail: "Internal error",
		},
	}

	for _, e := range testData {
		ne := New(e.ID, e.Detail, e.Code)

		if e.Error() != ne.Error() {
			t.Fatalf("Expected %s got %s", e.Error(), ne.Error())
		}

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.ID != e.ID {
			t.Fatalf("Expected %s got %s", e.ID, pe.ID)
		}

		if pe.Detail != e.Detail {
			t.Fatalf("Expected %s got %s", e.Detail, pe.Detail)
		}

		if pe.Code != e.Code {
			t.Fatalf("Expected %d got %d", e.Code, pe.Code)
		}
	}
}
