package xpubsub

import (
	"context"

	"gocloud.dev/pubsub"
)

// Middleware is a middleware function type.
type Middleware func(next Handler) Handler

// Middlewares is a slice of middlewares.
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

// Handle implements Handler interface.
func (c *ChainHandler) Handle(ctx context.Context, m *pubsub.Message) error {
	return c.chain.Handle(ctx, m)
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
