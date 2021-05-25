package xpubsub

import (
	"context"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/internal/uniqueid"
	"github.com/kucjac/cleango/xlog"
	"github.com/kucjac/cleango/xservice"
	"github.com/sirupsen/logrus"
	"gocloud.dev/pubsub"
)

// Compile time check if the Mux implements xservice.RunnerCloser interface.
var _ xservice.RunnerCloser = (*Mux)(nil)

type subscriber struct {
	sub         *pubsub.Subscription
	topic       string
	id          string
	maxHandlers int
}

// Mux is a subscriber router. It provides an easy interfaces for starting and listening on the subscriptions.
// Implements xservice.RunnerCloser interface.
type Mux struct {
	inline      bool
	maxHandlers int
	middlewares Middlewares
	routes      []route
	subRoutes   []subscriptionRoute
	parent      *Mux
	children    []*Mux
	ctx         context.Context
	cf          context.CancelFunc
	subscribers []subscriber
	running     bool
	gen         uniqueid.Generator
}

// NewMux creates a new mux that will register subscriptions using provided subscriber factory.
func NewMux() *Mux {
	return &Mux{ctx: context.Background(), maxHandlers: 10, gen: uniqueid.NextGenerator("pubsub:mux")}
}

// Run starts the router and begin on listening for given subscriptions.
func (m *Mux) Run() error {
	if m.running {
		return cgerrors.ErrInternal("mux is already running")
	}
	// Check if all routes are consistent and valid.
	mp := map[string]struct{}{}
	if err := m.checkRoutes(mp); err != nil {
		return err
	}

	m.ctx, m.cf = context.WithCancel(context.Background())
	if err := m.listenOnRoutes(); err != nil {
		m.cf()
		return err
	}
	if err := m.listenOnSubscriptions(); err != nil {
		m.cf()
		return err
	}
	m.running = true
	return nil
}

// Close the pubsub mux subscriptions.
func (m *Mux) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := m.close(ctx); err != nil {
		return err
	}
	for _, ch := range m.children {
		if err := ch.close(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Use provides middlewares to be used for all routes and children groups.
func (m *Mux) Use(middleware ...Middleware) {
	m.middlewares = append(m.middlewares, middleware...)
}

// With sets the middlewares for the resultant Mux that would be used exclusively in its and it's children context.
func (m *Mux) With(middlewares ...Middleware) *Mux {
	var mws []Middleware

	// Copy middlewares from parent mux.
	if m.inline {
		mws = make([]Middleware, len(m.middlewares))
		copy(mws, m.middlewares)
	}

	mws = append(mws, middlewares...)

	im := &Mux{inline: true, parent: m, middlewares: mws, maxHandlers: m.maxHandlers}
	m.children = append(m.children, im)
	return im
}

// WithMaxHandlers sets up the maximum concurrent handlers number for the resultant Mux that would be used exclusively in
// its and it's children context.
func (m *Mux) WithMaxHandlers(maxHandlers int) *Mux {
	var mws []Middleware
	if m.inline {
		mws = make([]Middleware, len(m.middlewares))
		copy(mws, m.middlewares)
	}

	im := &Mux{inline: true, parent: m, middlewares: mws, maxHandlers: maxHandlers}
	m.children = append(m.children, im)
	return im
}

// Subscribe registers topic subscriber that handles the message using provided handler with given options.
func (m *Mux) Subscribe(topic string, hf HandlerFunc) {
	m.routes = append(m.routes, route{
		topic:       topic,
		h:           hf,
		middlewares: m.middlewares,
		maxHandlers: m.maxHandlers,
	})
}

// Subscription registers subscription with specific handler.
// Optionally the topic name might be set for the logging and context passing purpose.
func (m *Mux) Subscription(sub *pubsub.Subscription, hf HandlerFunc, topic ...string) {
	sr := subscriptionRoute{
		sub:         sub,
		h:           hf,
		middlewares: m.middlewares,
		maxHandlers: m.maxHandlers,
	}
	if len(topic) > 0 {
		sr.topic = topic[0]
	}
	m.subRoutes = append(m.subRoutes, sr)
}

func (m *Mux) close(ctx context.Context) error {
	for _, s := range m.subscribers {
		fields := logrus.Fields{"id": s.id}
		if s.topic != "" {
			fields["topic"] = s.topic
		}
		xlog.WithFields(fields).Info("Closing subscription")
		if err := s.sub.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) listenOnRoutes() error {
	for _, r := range m.routes {

		sub, err := pubsub.OpenSubscription(m.ctx, r.topic)
		if err != nil {
			return err
		}

		// Create a subscription with it's unique id.
		sb := subscriber{
			sub:   sub,
			topic: r.topic,
			id:    m.gen.NextId(),
		}

		// Provide log fields for given subscription.
		logFields := logrus.Fields{
			"topic": r.topic,
			"id":    sb.id,
		}

		xlog.WithFields(logFields).Infof("listening for subscription")
		m.subscribers = append(m.subscribers, sb)

		go m.listenOnSubscriber(sub, sb.id, r.topic, r.maxHandlers, r.middlewares.Handler(r.h))
	}
	for _, ch := range m.children {
		if err := ch.listenOnRoutes(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) listenOnSubscriptions() error {
	for _, r := range m.subRoutes {
		// Create a subscription with it's unique id.
		sb := subscriber{
			sub:   r.sub,
			id:    m.gen.NextId(),
			topic: r.topic,
		}
		// Provide log fields for given subscription.
		logFields := logrus.Fields{
			"id": sb.id,
		}
		if r.topic != "" {
			logFields["topic"] = r.topic
		}

		xlog.WithFields(logFields).Info("listening for subscription")
		m.subscribers = append(m.subscribers, sb)

		go m.listenOnSubscriber(r.sub, sb.id, sb.topic, r.maxHandlers, r.middlewares.Handler(r.h))
	}
	for _, ch := range m.children {
		if err := ch.listenOnSubscriptions(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) checkRoutes(mp map[string]struct{}) error {
	for _, r := range m.routes {
		if r.topic == "" {
			return cgerrors.ErrInternal("no topic defined for one of the subscriber handlers")
		}
		_, ok := mp[r.topic]
		if ok {
			xlog.Warningf("topic: %s already has handler", r.topic)
		}
		if r.h == nil {
			return cgerrors.ErrInternalf("topic: %s handler not defined", r.topic)
		}
		mp[r.topic] = struct{}{}
	}
	for _, ch := range m.children {
		if err := ch.checkRoutes(mp); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) listenOnSubscriber(sb *pubsub.Subscription, id, topic string, maxHandlers int, handler Handler) {
	sem := make(chan struct{}, maxHandlers)
recvLoop:
	for {
		psMsg, err := sb.Receive(m.ctx)
		if err != nil {
			fields := logrus.Fields{
				"error": err,
			}
			if topic != "" {
				fields["topic"] = topic
			}
			xlog.WithFields(fields).Error("Receiving message failed")
			continue
		}

		// Wait if there are too many active handle goroutines and acquire the
		// semaphore. If the context is canceled, stop waiting and start shutting
		// down.
		select {
		case sem <- struct{}{}:
		case <-m.ctx.Done():
			break recvLoop
		}

		// Handle the message in a new goroutine.
		go func(msg *pubsub.Message, h Handler) {
			defer func() { <-sem }() // Release the semaphore.
			// An error should be
			ctx := context.WithValue(m.ctx, subIdCtxKey, id)
			if topic != "" {
				ctx = context.WithValue(ctx, subTopicCtxKey, topic)
			}
			_ = h.Handle(ctx, msg)
		}(psMsg, handler)
	}

	for n := 0; n < maxHandlers; n++ {
		sem <- struct{}{}
	}
}

type route struct {
	topic       string
	h           Handler
	middlewares Middlewares
	maxHandlers int
}

type subscriptionRoute struct {
	sub         *pubsub.Subscription
	h           Handler
	middlewares Middlewares
	maxHandlers int
	topic       string
}

type subscriptionCtx string

var (
	subTopicCtxKey = subscriptionCtx("sub:topic")
	subIdCtxKey    = subscriptionCtx("sub:id")
)
