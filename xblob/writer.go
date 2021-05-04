package xblob

import (
	"io"

	"gocloud.dev/blob"
)

//go:generate mockgen -destination=mockblob/writer_gen.go -package=mockblob . Writer

var _ Writer = (*blob.Writer)(nil)

type Writer interface {
	Write([]byte) (int, error)
	Close() error
	ReadFrom(r io.Reader) (int64, error)
}
