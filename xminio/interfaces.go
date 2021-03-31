package xminio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

// ObjectPutter is the interface used for uploading the object.
type ObjectPutter interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, input minio.PutObjectOptions) (minio.UploadInfo, error)
}

// ObjectGetter is the interface used for getting the object.
type ObjectGetter interface {
	GetObject(ctx context.Context, bucket, objectName string, options minio.GetObjectOptions) (*minio.Object, error)
}
