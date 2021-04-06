package xminio

import (
	"context"
	"io"
	"os"

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

// FileObject wraps input file to implement Object interface.
func FileObject(f *os.File) Object {
	return &fileObject{f: f}
}

var _ Object = (*fileObject)(nil)

type fileObject struct {
	f *os.File
}

func (f *fileObject) Read(b []byte) (int, error) {
	return f.f.Read(b)
}

func (f *fileObject) Stat() (minio.ObjectInfo, error) {
	st, err := f.f.Stat()
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	return minio.ObjectInfo{
		Key:          st.Name(),
		LastModified: st.ModTime(),
		Size:         st.Size(),
	}, nil
}

func (f *fileObject) ReadAt(b []byte, offset int64) (int, error) {
	return f.f.ReadAt(b, offset)
}

func (f *fileObject) Seek(offset int64, whence int) (int64, error) {
	return f.f.Seek(offset, whence)
}

func (f *fileObject) Close() error {
	return f.f.Close()
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
