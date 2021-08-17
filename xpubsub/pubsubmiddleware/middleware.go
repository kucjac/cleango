package pubsubmiddleware

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/kucjac/cleango/xpubsub"
	"github.com/sirupsen/logrus"
	"gocloud.dev/pubsub"
	"google.golang.org/grpc/metadata"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"

	"github.com/kucjac/cleango/internal/uniqueid"
)

// Key to use when setting the message ID.
type ctxKeyMessageID int

var subIdGen = uniqueid.NextGenerator("subscription")

// MessageIDKey is the key that holds the unique request ID in a request context.
const MessageIDKey ctxKeyMessageID = 0

// CtxSubject gets the subscription subject from the given context.
func CtxSubject(ctx context.Context) string {
	t, _ := ctx.Value(xpubsub.SubscriptionSubjectCtxKey).(string)
	return t
}

// CtxSubscriptionID gets the subscription id from the given context.
func CtxSubscriptionID(ctx context.Context) string {
	id, _ := ctx.Value(xpubsub.SubscriptionIdCtxKey).(string)
	return id
}

// Acker is a middleware that checks if the subsequent handlers returns an error and on success Ackes them.
// In case of failure it Nacks the message if it is possible.
func Acker(h xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, m *pubsub.Message) error {
		if err := h.Handle(ctx, m); err != nil {
			if m.Nackable() {
				m.Nack()
			} else {
				// If the implementation doesn't allow us to Nack the message it must be Acknowledged.
				m.Ack()
			}
			return err
		}
		m.Ack()
		return nil
	})
}

// Recoverer recovers from any panic in the handler and appends RecoveredPanicError with the stacktrace
// to any error returned from the handler.
func Recoverer(h xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, m *pubsub.Message) error {
		defer func() {
			if r := recover(); r != nil {
				xlog.Errorf("panic occurred: %#v, stack: \n%s", r, string(debug.Stack()))
			}
		}()
		if err := h.Handle(ctx, m); err != nil {
			return err
		}
		return nil
	})
}

// Logger is a middleware function that is used for trace logging incoming message on subscriptions.
func Logger(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, m *pubsub.Message) error {
		fields := logrus.Fields{
			"subscriptionId": CtxSubscriptionID(ctx),
			"topic":          CtxSubject(ctx),
		}
		reqID := GetMessageID(ctx)
		if reqID != "" {
			fields["messageId"] = reqID
		}
		ts := time.Now()

		var (
			msg  string
			code cgerrors.ErrorCode
		)
		err := next.Handle(ctx, m)
		if err != nil {
			if e, ok := err.(*cgerrors.Error); ok {
				code = e.Code
				msg = e.Detail
			} else {
				code = cgerrors.Code(err)
				msg = err.Error()
			}
		}
		fields["code"] = code
		if msg != "" {
			fields["detail"] = msg
		}
		xlog.WithContext(ctx).
			WithFields(fields).
			Tracef("message handled in %s", time.Since(ts))
		return err
	})
}

// MessageID is a middleware function that generates new request id and puts it into message context,
// A message ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
// The concept and implementation of this request id is based on the brilliant golang library: github.com/go-chi/chi.
func MessageID(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, m *pubsub.Message) error {
		if ctx == nil {
			ctx = context.Background()
		}
		messageId := subIdGen.NextId()
		ctx = context.WithValue(ctx, MessageIDKey, messageId)

		return next.Handle(ctx, m)
	})
}

// ContextMetadata sets up the metadata from the message in the context.
func ContextMetadata(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, m *pubsub.Message) error {
		if ctx == nil {
			ctx = context.Background()
		}
		md := m.Metadata
		if md == nil {
			md = map[string]string{}
		}
		ctx = metadata.NewIncomingContext(ctx, metadata.New(md))
		return next.Handle(ctx, m)
	})
}

// GetMessageID returns a message ID from the given context if one is present.
// Returns the empty string if a message ID cannot be found.
func GetMessageID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(MessageIDKey).(string); ok {
		return reqID
	}
	return ""
}
