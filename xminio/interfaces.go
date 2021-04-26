package xminio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

//go:generate mockgen -destination=xminiomock/object_gen.go -package=xminiomock . ObjectPutterGetter,Object

// ObjectPutter is the interface used for uploading the object.
type ObjectPutter interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, input minio.PutObjectOptions) (minio.UploadInfo, error)
}

// ObjectGetter is the interface used for getting the object.
type ObjectGetter interface {
	GetObject(ctx context.Context, bucket, objectName string, options minio.GetObjectOptions) (Object, error)
}

// ObjectPutterGetter is the interface used for getting and putting objects into minio s3 like storage.
type ObjectPutterGetter interface {
	ObjectPutter
	ObjectGetter
}

// Object is an interface that has exact methods as minio.Object.
type Object interface {
	Read(b []byte) (int, error)
	Stat() (minio.ObjectInfo, error)
	ReadAt(b []byte, offset int64) (int, error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

var _ ObjectPutterGetter = (*wrapped)(nil)

// WrapClient wraps *minio.Client to implement ObjectPutterGetter interface.
func WrapClient(c *minio.Client) ObjectPutterGetter {
	return &wrapped{Client: c}
}

type wrapped struct {
	Client *minio.Client
}

func (w *wrapped) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, options minio.PutObjectOptions) (minio.UploadInfo, error) {
	return w.Client.PutObject(ctx, bucketName, objectName, reader, objectSize, options)
}

func (w *wrapped) GetObject(ctx context.Context, bucket, objectName string, options minio.GetObjectOptions) (Object, error) {
	return w.Client.GetObject(ctx, bucket, objectName, options)
}
