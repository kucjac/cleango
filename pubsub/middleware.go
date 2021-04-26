package pubsub

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

var (
	prefix string
	reqId  uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// CtxTopic gets the subscription topic from the given context.
func CtxTopic(ctx context.Context) string {
	t, _ := ctx.Value(subTopicCtxKey).(string)
	return t
}

// CtxSubscriptionID gets the subscription id from the given context.
func CtxSubscriptionID(ctx context.Context) string {
	id, _ := ctx.Value(subIdCtxKey).(string)
	return id
}

// Middleware is a middleware function type.
type Middleware func(next Handler) Handler

type Middlewares []Middleware

func (mws Middlewares) Handler(h Handler) Handler {
	return &ChainHandler{Middlewares: mws, Endpoint: h, chain: chain(mws, h)}
}

func (mws Middlewares) HandlerFunc(h HandlerFunc) Handler {
	return &ChainHandler{Middlewares: mws, Endpoint: h, chain: chain(mws, h)}
}

// ChainHandler is a http.Handler with support for handler composition and
// execution.
type ChainHandler struct {
	Middlewares Middlewares
	Endpoint    Handler
	chain       Handler
}

func (c *ChainHandler) Handle(m *message.Message) {
	c.chain.Handle(m)
}

// chain builds a http.Handler composed of an inline middleware stack and endpoint
// handler in the order they are passed.
func chain(middlewares Middlewares, endpoint Handler) Handler {
	// Return ahead of time if there aren't any middlewares for the chain
	if len(middlewares) == 0 {
		return endpoint
	}

	// Wrap the end handler with the middleware chain
	h := middlewares[len(middlewares)-1](endpoint)
	for i := len(middlewares) - 2; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

// Recoverer recovers from any panic in the handler and appends RecoveredPanicError with the stacktrace
// to any error returned from the handler.
func Recoverer(h Handler) Handler {
	return HandlerFunc(func(m *message.Message) {
		defer func() {
			if r := recover(); r != nil {
				xlog.Errorf("panic occurred: %#v, stack: \n%s", r, string(debug.Stack()))
			}
		}()
		h.Handle(m)
	})
}

// Logger is a middleware function that is used for trace logging incoming message on subscriptions.
func Logger(next Handler) Handler {
	return HandlerFunc(func(m *message.Message) {
		ctx := m.Context()
		fields := logrus.Fields{
			"messageId":      m.UUID,
			"subscriptionId": CtxSubscriptionID(ctx),
			"topic":          CtxTopic(ctx),
		}
		ts := time.Now()
		next.Handle(m)
		xlog.WithFields(fields).
			Tracef("message handled in %s", time.Since(ts))
	})
}

// RequestID is a middleware function that generates new request id and puts it into message context,
// A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
// The concept and implementation of this request id is based on the brilliant golang library: github.com/go-chi/chi.
func RequestID(next Handler) Handler {
	return HandlerFunc(func(m *message.Message) {
		ctx := m.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		thisID := atomic.AddUint64(&reqId, 1)
		requestID := fmt.Sprintf("%s-%06d", prefix, thisID)
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		m.SetContext(ctx)

		next.Handle(m)
	})
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
