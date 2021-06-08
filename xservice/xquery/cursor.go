package xquery

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/gob"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"
)

// CursorType is the enumerator that defines the type of the cursor.
type CursorType uint8

// Cursors enumerated types. There are 5 defined cursor types:
//	- This 	-
//	- Prev	-
//	- Next	-
//	- First	-
//	- Last	-
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

// EncodeCursor encodes the cursor into a string.
func EncodeCursor(m interface{}) (string, error) {
	var (
		data []byte
		err  error
	)
	switch mt := m.(type) {
	case encoding.BinaryMarshaler:
		data, err = mt.MarshalBinary()
		if err != nil {
			return "", cgerrors.ErrInternal("marshaling curo")
		}
	case nil:
		return "", cgerrors.ErrInternal("provided nil cursor to encode")
	default:
		buf := bytes.Buffer{}
		e := gob.NewEncoder(&buf)
		if err = e.Encode(m); err != nil {
			return "", cgerrors.ErrInternal("encoding")
		}
		data = buf.Bytes()
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeCursor decodes the cursor entry input.
func DecodeCursor(in string, m interface{}) error {
	bt, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return cgerrors.ErrInvalidArgument("invalid cursor")
	}
	if len(bt) == 0 {
		return cgerrors.ErrInvalidArgument("empty cursor")
	}

	switch mt := m.(type) {
	case encoding.BinaryUnmarshaler:
		if err = mt.UnmarshalBinary(bt); err != nil {
			return err
		}
	case nil:
		return cgerrors.ErrInternal("provided nil cursor to decode")
	default:
		d := gob.NewDecoder(bytes.NewReader(bt))
		if err = d.Decode(m); err != nil {
			xlog.WithField("err", err).
				Debug("decoding cursor failed")
			return cgerrors.ErrInvalidArgument("invalid cursor")
		}
	}
	return nil
}
