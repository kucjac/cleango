package xminio

import (
	"bytes"
	"time"

	"github.com/minio/minio-go/v7"
)

var _ Object = (*bytesObject)(nil)

type bytesObject struct {
	lastModified time.Time
	r            *bytes.Reader
	size         int
}

func NewBytesObject(data []byte, lastModified time.Time) Object {
	return &bytesObject{
		lastModified: lastModified,
		r:            bytes.NewReader(data),
		size:         len(data),
	}
}

func (bo *bytesObject) Read(b []byte) (int, error) {
	return bo.r.Read(b)
}

func (bo *bytesObject) Stat() (minio.ObjectInfo, error) {
	return minio.ObjectInfo{LastModified: bo.lastModified, Size: int64(bo.size)}, nil
}

func (bo *bytesObject) ReadAt(b []byte, offset int64) (int, error) {
	return bo.r.ReadAt(b, offset)
}

func (bo *bytesObject) Seek(offset int64, whence int) (int64, error) {
	return bo.r.Seek(offset, whence)
}

func (bo *bytesObject) Close() error {
	return nil
}
