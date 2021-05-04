package xblob

import (
	"context"

	"gocloud.dev/blob"
)

//go:generate mockgen -destination=mockblob/bucket_gen.go -package=mockblob . Bucket

// Bucket is the wrapped interface with small extension that returns Reader and Writer interfaces instead of specific implementations.
type Bucket interface {
	Bucket() *blob.Bucket
	As(i interface{}) bool
	ErrorAs(err error, i interface{}) bool
	ReadAll(ctx context.Context, key string) (_ []byte, err error)
	List(opts *blob.ListOptions) *blob.ListIterator
	ListPage(ctx context.Context, pageToken []byte, pageSize int, opts *blob.ListOptions) (retval []*blob.ListObject, nextPageToken []byte, err error)
	IsAccessible(ctx context.Context) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
	Attributes(ctx context.Context, key string) (_ *blob.Attributes, err error)
	NewReader(ctx context.Context, key string, opts *blob.ReaderOptions) (Reader, error)
	NewRangeReader(ctx context.Context, key string, offset, length int64, opts *blob.ReaderOptions) (_ Reader, err error)
	WriteAll(ctx context.Context, key string, p []byte, opts *blob.WriterOptions) (err error)
	NewWriter(ctx context.Context, key string, opts *blob.WriterOptions) (_ Writer, err error)
	Copy(ctx context.Context, dstKey, srcKey string, opts *blob.CopyOptions) (err error)
	Delete(ctx context.Context, key string) (err error)
	SignedURL(ctx context.Context, key string, opts *blob.SignedURLOptions) (string, error)
	Close() error
}

// Wrap the input blob to implement Bucket interface.
func Wrap(b *blob.Bucket) Bucket {
	return &bucket{b: b}
}

// bucket is a wrapper over blob.Bucket that implements Bucket interface.
type bucket struct {
	b *blob.Bucket
}

func (b *bucket) Bucket() *blob.Bucket {
	return b.b
}

func (b *bucket) As(i interface{}) bool {
	return b.b.As(i)
}

func (b *bucket) ErrorAs(err error, i interface{}) bool {
	return b.b.ErrorAs(err, i)
}

func (b *bucket) ReadAll(ctx context.Context, key string) (_ []byte, err error) {
	return b.b.ReadAll(ctx, key)
}

func (b *bucket) List(opts *blob.ListOptions) *blob.ListIterator {
	return b.b.List(opts)
}

func (b *bucket) ListPage(ctx context.Context, pageToken []byte, pageSize int, opts *blob.ListOptions) (retval []*blob.ListObject, nextPageToken []byte, err error) {
	return b.b.ListPage(ctx, pageToken, pageSize, opts)
}

func (b *bucket) IsAccessible(ctx context.Context) (bool, error) {
	return b.b.IsAccessible(ctx)
}

func (b *bucket) Exists(ctx context.Context, key string) (bool, error) {
	return b.b.Exists(ctx, key)
}

func (b *bucket) Attributes(ctx context.Context, key string) (_ *blob.Attributes, err error) {
	return b.b.Attributes(ctx, key)
}

func (b *bucket) NewReader(ctx context.Context, key string, opts *blob.ReaderOptions) (Reader, error) {
	return b.b.NewReader(ctx, key, opts)
}

func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *blob.ReaderOptions) (_ Reader, err error) {
	return b.b.NewRangeReader(ctx, key, offset, length, opts)
}

func (b *bucket) WriteAll(ctx context.Context, key string, p []byte, opts *blob.WriterOptions) (err error) {
	return b.b.WriteAll(ctx, key, p, opts)
}

func (b *bucket) NewWriter(ctx context.Context, key string, opts *blob.WriterOptions) (_ Writer, err error) {
	return b.b.NewWriter(ctx, key, opts)
}

func (b *bucket) Copy(ctx context.Context, dstKey, srcKey string, opts *blob.CopyOptions) (err error) {
	return b.b.Copy(ctx, dstKey, srcKey, opts)
}

func (b *bucket) Delete(ctx context.Context, key string) (err error) {
	return b.b.Delete(ctx, key)
}

func (b *bucket) SignedURL(ctx context.Context, key string, opts *blob.SignedURLOptions) (string, error) {
	return b.b.SignedURL(ctx, key, opts)
}

func (b *bucket) Close() error {
	return b.b.Close()
}
