package xhttp

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes the JSON output in given response writer.
func WriteJSON(rw http.ResponseWriter, in interface{}, options ...ResponseOption) {
	o := &ResponseOptions{Headers: map[string]string{"Content-Type": "application/json"}}
	for _, option := range options {
		option(o)
	}

	data, err := json.Marshal(in)
	if err != nil {
		writeUndefinedError(err, rw, nil)
		return
	}

	writeHeaders(rw, o)
	c := http.StatusOK
	if o.Status != 0 {
		c = o.Status
	}
	rw.WriteHeader(c)
	rw.Write(data)
}
