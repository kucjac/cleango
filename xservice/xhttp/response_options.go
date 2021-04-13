package xhttp

// ResponseWithStatus returns an error option that sets given response status.
func ResponseWithStatus(status int) ResponseOption {
	return func(e *ResponseOptions) {
		e.Status = status
	}
}

// ResponseWithHeader is an option to the http response that sets given header with corresponding value
func ResponseWithHeader(key, value string) ResponseOption {
	return func(e *ResponseOptions) {
		if e.Headers == nil {
			e.Headers = make(map[string]string)
		}
		e.Headers[key] = value
	}
}

// ResponseWithContentType is an option to the http response that sets provided content type.
func ResponseWithContentType(contentType string) ResponseOption {
	return func(e *ResponseOptions) {
		if e.Headers == nil {
			e.Headers = make(map[string]string)
		}
		e.Headers["Content-Type"] = contentType
	}
}

// ResponseOption is an option for the error function.
type ResponseOption func(e *ResponseOptions)

// ResponseOptions are the options used to customize HTTP response.
type ResponseOptions struct {
	Status  int
	Headers map[string]string
}
