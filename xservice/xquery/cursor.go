package xquery

import (
	"bytes"
	"encoding"
	"encoding/base64"

	"github.com/kucjac/cleango/cgerrors"
)

// CursorType is the enumerator that defines the type of the cursor.
type CursorType uint8

// Cursors enumerated types.
const (
	CursorTypeUndefined CursorType = iota
	CursorTypeThis
	CursorTypePrev
	CursorTypeNext
	CursorTypeFirst
	CursorTypeLast
)

// IsValid checks if the provided cursor type is valid.
func (c CursorType) IsValid() bool {
	return c > CursorTypeUndefined && c <= CursorTypeLast
}

// CursorEntry is the entry used by the cursor.
type CursorEntry struct {
	Type  CursorType
	Value string
}

// IsNull checks if the cursor entry is null.
func (c *CursorEntry) IsNull() bool {
	return !c.Type.IsValid() && c.Value == ""
}

func (c CursorEntry) Encode() string {
	b := bytes.Buffer{}
	b.WriteByte(uint8(c.Type))
	b.Write([]byte(c.Value))
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// Decode the cursor entry.
func (c *CursorEntry) Decode(in string) error {
	bt, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	if len(bt) == 0 {
		return cgerrors.ErrInvalidArgument("empty cursor")
	}

	c.Type = CursorType(bt[0])
	if c.Type == CursorTypeUndefined || c.Type > CursorTypeLast {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	if len(bt) > 1 {
		c.Value = string(bt[1:])
	}
	return nil
}

// CursorEntryMarshaler is an interace used to marshal custom cursor entries.
type CursorEntryMarshaler interface {
	encoding.BinaryMarshaler
	GetType() CursorType
}

// CursorEntryUnmarshaler is the interface used to unmarshal the custom cursor entry.
type CursorEntryUnmarshaler interface {
	encoding.BinaryUnmarshaler
	SetType(t CursorType)
}

// EncodeCursorEntry encodes the cursor into a string.
func EncodeCursorEntry(m CursorEntryMarshaler) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteByte(byte(m.GetType()))

	data, err := m.MarshalBinary()
	if err != nil {
		return "", err
	}
	buf.Write(data)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DecodeCursorEntry decodes the cursor entry input.
func DecodeCursorEntry(in string, u CursorEntryUnmarshaler) error {
	bt, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	if len(bt) == 0 {
		return cgerrors.ErrInvalidArgument("empty cursor")
	}

	tp := CursorType(bt[0])
	if tp == CursorTypeUndefined || tp > CursorTypeLast {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	u.SetType(tp)
	if err = u.UnmarshalBinary(bt[1:]); err != nil {
		return err
	}
	return nil
}
