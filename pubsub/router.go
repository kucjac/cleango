package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kucjac/cleango/errors"
	"github.com/kucjac/cleango/xlog"
	"github.com/kucjac/cleango/xservice"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Compile time check if the Mux implements xservice.RunnerCloser interface.
var _ xservice.RunnerCloser = (*Mux)(nil)

type subscriber struct {
	sub   Subscriber
	topic string
	id    string
}

// Mux is a subscriber router. It provides an easy interfaces for starting and listening on the subscriptions.
// Implements xservice.RunnerCloser interface.
type Mux struct {
	sf          SubscriberFactory
	inline      bool
	middlewares Middlewares
	routes      []route
	parent      *Mux
	children    []*Mux
	ctx         context.Context
	cf          context.CancelFunc
	subscribers []subscriber
	running     bool
}

// NewMux creates a new mux that will register subscriptions using provided subscriber factory.
func NewMux(sf SubscriberFactory) *Mux {
	return &Mux{sf: sf}
}

// Run starts the router and begin on listening for given subscriptions.
func (m *Mux) Run() error {
	if m.running {
		return errors.ErrInternal("mux is already running")
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
	m.running = true
	return nil
}

func (m *Mux) Close() error {
	if err := m.close(); err != nil {
		return err
	}
	for _, ch := range m.children {
		if err := ch.close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) close() error {
	for _, s := range m.subscribers {
		fields := logrus.Fields{
			"id":    s.id,
			"topic": s.topic,
		}
		xlog.WithFields(fields).Infof("Closing subscription: %s", s.topic)
		if err := s.sub.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) listenOnRoutes() error {
	for _, r := range m.routes {
		s, err := m.sf.NewSubscriber(r.options...)
		if err != nil {
			return err
		}

		// Create a subscription with it's unique id.
		sb := subscriber{
			sub:   s,
			topic: r.topic,
			id:    uuid.NewV4().String(),
		}

		// Provide log fields for given subscription.
		logFields := logrus.Fields{
			"topic": r.topic,
			"id":    sb.id,
		}
		if len(r.options) > 0 {
			so := &SubscriptionOptions{}
			for _, o := range r.options {
				o(so)
			}
			if so.DurableName != "" {
				logFields["durable_name"] = so.DurableName
			}
			if so.QueueGroup != "" {
				logFields["queue_group"] = so.QueueGroup
			}
		}
		xlog.WithFields(logFields).Infof("listening at the topic: %s", r.topic)
		mc, err := s.Subscribe(m.ctx, r.topic)
		if err != nil {
			return err
		}
		m.subscribers = append(m.subscribers, sb)

		go func(mc <-chan *message.Message, id, topic string, handler Handler) {
			for m := range mc {
				ctx := context.WithValue(m.Context(), subIdCtxKey, id)
				ctx = context.WithValue(ctx, subTopicCtxKey, topic)
				handler.Handle(m)
			}
		}(mc, sb.id, r.topic, r.middlewares.Handler(r.h))
	}
	for _, ch := range m.children {
		if err := ch.listenOnRoutes(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mux) checkRoutes(mp map[string]struct{}) error {
	for _, r := range m.routes {
		if r.topic == "" {
			return errors.ErrInternal("no topic defined for one of the subscriber handlers")
		}
		_, ok := mp[r.topic]
		if ok {
			xlog.Warningf("topic: %s already has handler", r.topic)
		}
		if r.h == nil {
			return errors.ErrInternalf("topic: %s handler not defined", r.topic)
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

// Use provides middlewares to be used for all routes and children groups.
func (m *Mux) Use(middleware ...Middleware) {
	m.middlewares = append(m.middlewares, middleware...)
}

// With sets the middlewares for the resultant Mux that would be used exclusively in its and it's childrens context.
func (m *Mux) With(middlewares ...Middleware) *Mux {
	var mws []Middleware

	// Copy middlewares from parent mux.
	if m.inline {
		mws = make([]Middleware, len(m.middlewares))
		copy(mws, m.middlewares)
	}

	mws = append(mws, middlewares...)

	im := &Mux{inline: true, parent: m, middlewares: mws}
	m.children = append(m.children, im)
	return im
}

// Subscribe registers topic subscriber that handles the message using provided handler with given options.
func (m *Mux) Subscribe(topic string, handler Handler, options ...SubscriptionOption) {
	m.routes = append(m.routes, route{
		topic:       topic,
		h:           handler,
		middlewares: m.middlewares,
		options:     options,
	})
}

type route struct {
	topic       string
	h           Handler
	middlewares Middlewares
	options     []SubscriptionOption
}

type subscriptionCtx string

var (
	subTopicCtxKey = subscriptionCtx("sub:topic")
	subIdCtxKey    = subscriptionCtx("sub:id")
)
