package xhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/kucjac/cleango/cgerrors"
)

// WriteErrJSON writes the JSON error into given response writer.
func WriteErrJSON(rw http.ResponseWriter, err error, options ...ResponseOption) {
	o := &ResponseOptions{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	for _, option := range options {
		option(o)
	}

	var e *cgerrors.Error
	if errors.As(err, &e) {
		writeDefinedError(rw, e, o)
		return
	}
	writeUndefinedError(err, rw, o)

}

func writeDefinedError(rw http.ResponseWriter, e *cgerrors.Error, o *ResponseOptions) {
	meta := e.Meta
	if !o.Debug && e.Meta != nil {
		_, ok := e.Meta[cgerrors.MetaKeyWrapped]
		if ok {
			temp := map[string]string{}
			for k, v := range e.Meta {
				if k == cgerrors.MetaKeyWrapped {
					continue
				}
				temp[k] = v
			}
			e.Meta = temp
		}
	}
	data, err := json.Marshal(e)
	if err != nil {
		writeUndefinedError(e, rw, o)
		return
	}
	e.Meta = meta

	c := o.Status
	if c == 0 {
		switch codes.Code(e.Code) {
		case codes.Internal:
			c = http.StatusInternalServerError
		case codes.InvalidArgument:
			c = http.StatusBadRequest
		case codes.PermissionDenied:
			c = http.StatusForbidden
		case codes.Unauthenticated:
			c = http.StatusUnauthorized
		case codes.AlreadyExists:
			c = http.StatusConflict
		case codes.NotFound:
			c = http.StatusNotFound
		case codes.Unavailable:
			c = http.StatusServiceUnavailable
		case codes.DeadlineExceeded:
			c = http.StatusGatewayTimeout
		default:
			c = http.StatusInternalServerError
		}
	}
	writeHeaders(rw, o)
	rw.WriteHeader(c)
	rw.Write(data)
}

func writeUndefinedError(err error, rw http.ResponseWriter, o *ResponseOptions) {
	type jsonError struct {
		Message string `json:"message"`
	}
	data, er := json.Marshal(jsonError{Message: err.Error()})
	if er != nil {
		data = []byte(`{"message":"internal server error"}`)
	}
	c := http.StatusInternalServerError
	if o != nil && o.Status != 0 {
		c = o.Status
	}
	writeHeaders(rw, o)
	rw.WriteHeader(c)
	rw.Write(data)
}

func writeHeaders(rw http.ResponseWriter, o *ResponseOptions) {
	if o == nil || len(o.Headers) == 0 {
		return
	}
	header := rw.Header()
	for k, v := range o.Headers {
		header.Set(k, v)
	}
}
