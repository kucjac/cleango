package xminio

import (
	"os"

	"github.com/minio/minio-go/v7"
)

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
