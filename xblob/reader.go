package xblob

import (
	"io"
	"time"

	"gocloud.dev/blob"
)

//go:generate mockgen -destination=mockblob/reader_gen.go -package=mockblob . Reader

var _ Reader = (*blob.Reader)(nil)

// Reader is the interface implementation of the gocloud.dev/blob.Reader
// Using an interface allows to mock up the reader.
type Reader interface {
	Read([]byte) (int, error)
	Close() error
	ContentType() string
	ModTime() time.Time
	Size() int64
	WriteTo(w io.Writer) (int64, error)
}
