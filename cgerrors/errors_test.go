package cgerrors

import (
	er "errors"
	"testing"
)

func TestFromError(t *testing.T) {
	var err error
	err = ErrNotFoundf("%s", "example")
	merr := FromError(err)
	if merr.Code != ErrorCode_NotFound {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
	err = er.New(err.Error())
	merr = FromError(err)
	if merr.Code != ErrorCode_NotFound {
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
			ID:     "test",
			Code:   ErrorCode_Internal,
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
