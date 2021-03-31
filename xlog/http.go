package xlog

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// HTTPRequest for logging HTTP requests. Only contains semantics
// defined by the HTTP specification. Product-specific logging
// information MUST be defined in a separate message.
type HTTPRequest struct {
	// The request method. Examples: `"GET"`, `"HEAD"`, `"PUT"`, `"POST"`.
	RequestMethod string `json:"requestMethod,omitempty"`
	// The scheme (http, https), the host name, the path and the query
	// portion of the URL that was requested.
	// Example: `"http://example.com/some/info?color=red"`.
	RequestURL string `json:"requestUrl,omitempty"`
	// The size of the HTTP request message in bytes, including the request
	// headers and the request body.
	RequestSize int64 `json:"requestSize,omitempty"`
	// The response code indicating the status of response.
	// Examples: 200, 404.
	Status int32 `json:"status,omitempty"`
	// The size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize int64 `json:"responseSize,omitempty"`
	// The user agent sent by the client. Example:
	// `"Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)"`.
	UserAgent string `json:"userAgent,omitempty"`
	// The IP address (IPv4 or IPv6) of the client that issued the HTTP
	// request. Examples: `"192.168.1.1"`, `"FE80::0202:B3FF:FE1E:8329"`.
	RemoteIP string `json:"remoteIp,omitempty"`
	// The IP address (IPv4 or IPv6) of the origin server that the request was
	// sent to.
	ServerIP string `json:"serverIp,omitempty"`
	// The referer URL of the request, as defined in
	// [HTTP/1.1 Header Field Definitions](http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html).
	Referer string `json:"referer,omitempty"`
	// The request processing latency on the server, from the time the request was
	// received until the response was sent.
	Latency string `json:"latency,omitempty"`
	// Whether or not a cache lookup was attempted.
	CacheLookup bool `json:"cacheLookup,omitempty"`
	// Whether or not an entity was served from cache
	// (with or without validation).
	CacheHit bool `json:"cacheHit,omitempty"`
	// Whether or not the response was validated with the origin server before
	// being served from cache. This field is only meaningful if `cache_hit` is
	// True.
	CacheValidatedWithOriginServer bool `json:"cacheValidatedWithOriginServer,omitempty"`
	// The number of HTTP response bytes inserted into cache. Set only when a
	// cache fill was attempted.
	CacheFillBytes int64 `json:"cacheFillBytes,omitempty"`
	// Protocol used for the request. Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `json:"protocol,omitempty"`
}

// HTTPRequestKey key used in fields
// to log http request data
const HTTPRequestKey = "http:request"



// ReqWrap will put xlog into request context
func ReqWrap(req *http.Request, entry *logrus.Entry) *http.Request {
	return req.WithContext(CtxPut(req.Context(), entry))
}

// Req will pull xlog from request context
func Req(req *http.Request) *logrus.Entry {
	return Ctx(req.Context())
}

// PrintfHTTPReq will log request on specific level
func PrintfHTTPReq(ctx context.Context, req HTTPRequest) {
	Ctx(ctx).WithField(HTTPRequestKey, req).Printf("Request %s %d %s", req.RequestMethod, req.Status, req.RequestURL)
}
