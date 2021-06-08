package xquery

import (
	"bytes"

	"github.com/kucjac/cleango/cgerrors"
)

// CursorEntry is the entry used by the cursor.
type CursorEntry struct {
	Type  CursorType
	Value string
}

// IsNull checks if the cursor entry is null.
func (c *CursorEntry) IsNull() bool {
	return !c.Type.IsValid() && c.Value == ""
}

// MarshalBinary implements encoding.BinaryMarshaler interface.
func (c CursorEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte(uint8(c.Type))
	b.Write([]byte(c.Value))
	return b.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface.
func (c *CursorEntry) UnmarshalBinary(in []byte) error {
	if len(in) == 0 {
		return cgerrors.ErrInvalidArgument("invalid cursor data")
	}
	c.Type = CursorType(in[0])
	if c.Type == CursorTypeUndefined || c.Type > CursorTypeLast {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	if len(in) > 1 {
		c.Value = string(in[1:])
	}
	return nil
}
