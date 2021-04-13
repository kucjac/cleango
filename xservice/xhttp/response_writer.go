package xhttp

import (
	"bytes"
	"net/http"
)

var _ http.ResponseWriter = &ResponseWriter{}

// ResponseWriter is a wrapper implementation for the http.ResponseWriter. It allows to have buffered writer and status.
// It needs to be used in a middleware so that it could manipulate the output with some compressor or react differently on panic.
// Look at middleware.ResponseWriter how it is used.
type ResponseWriter struct {
	Status  int
	Buffer  *bytes.Buffer
	Wrapped http.ResponseWriter
}

// Header implements http.ResponseWriter.
func (r *ResponseWriter) Header() http.Header {
	return r.Wrapped.Header()
}

// Write implements io.Writer.
func (r *ResponseWriter) Write(bytes []byte) (int, error) {
	return r.Buffer.Write(bytes)
}

// WriteHeader implements http.ResponseWriter.
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Status = statusCode
}

// WrapResponseWriter wraps given response writer and creates an implementation that holds temporary buffer and status.
func WrapResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		Buffer:  bytes.NewBuffer(nil),
		Wrapped: rw,
	}
}
