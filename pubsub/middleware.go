package pubsub

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kucjac/cleango/xlog"
	"github.com/sirupsen/logrus"
)

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
