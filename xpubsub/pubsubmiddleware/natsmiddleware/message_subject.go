package natsmiddleware

import (
	"context"

	"github.com/nats-io/nats.go"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/natspubsub"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xpubsub"
)

// JetStreamMetadataCtxKey is the key to the JetStream Message metadata.
type JetStreamMetadataCtxKey struct{}

// Compile time check if the ExtractJetStreamMetadata function is xpubsub.Middleware.
var _ xpubsub.Middleware = ExtractJetStreamMetadata

// ExtractJetStreamMetadata is a middleware that extracts JetStream message metadata and stores in the context
// with the key JetStreamMetadataCtxKey.
func ExtractJetStreamMetadata(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, msg *pubsub.Message) error {
		var natsMsg *nats.Msg
		if !msg.As(&natsMsg) {
			return cgerrors.ErrInternal("pubsub.Message.As returned failed").
				WithMeta("msgId", msg.LoggableID)
		}
		// Get JetStream metadata.
		meta, err := natsMsg.Metadata()
		if err != nil {
			return err
		}

		// Set up the timestamp in the context.
		ctx = context.WithValue(ctx, JetStreamMetadataCtxKey{}, meta)

		// Handle the next handlers.
		if err := next.Handle(ctx, msg); err != nil {
			return err
		}
		return nil
	})
}

// Compile time check if the ExtractJetStreamMsgTimestamp function is xpubsub.Middleware.
var _ xpubsub.Middleware = ExtractJetStreamMsgTimestamp

// ExtractJetStreamMsgTimestamp extracts the timestamp from the JetStream Message metadata and stores it in the context.
func ExtractJetStreamMsgTimestamp(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, msg *pubsub.Message) error {
		meta, ok := ctx.Value(JetStreamMetadataCtxKey{}).(*nats.MsgMetadata)
		if !ok {
			var natsMsg *nats.Msg
			if !msg.As(&natsMsg) {
				return cgerrors.ErrInternal("pubsub.Message.As returned failed").
					WithMeta("msgId", msg.LoggableID)
			}
			// Get JetStream metadata.
			var err error
			meta, err = natsMsg.Metadata()
			if err != nil {
				return err
			}
		}

		// Set up the timestamp in the context.
		ctx = xpubsub.CtxSetMessageTimestamp(ctx, meta.Timestamp)

		// Handle the next handlers.
		if err := next.Handle(ctx, msg); err != nil {
			return err
		}
		return nil
	})
}

// Compile time check if the ExtractJetStreamMsgSequence function is xpubsub.Middleware.
var _ xpubsub.Middleware = ExtractJetStreamMsgSequence

// ExtractJetStreamMsgSequence extracts the timestamp from the JetStream Message metadata and stores it in the context.
func ExtractJetStreamMsgSequence(next xpubsub.Handler) xpubsub.Handler {
	return xpubsub.HandlerFunc(func(ctx context.Context, msg *pubsub.Message) error {
		meta, ok := ctx.Value(JetStreamMetadataCtxKey{}).(*nats.MsgMetadata)
		if !ok {
			var natsMsg *nats.Msg
			if !msg.As(&natsMsg) {
				return cgerrors.ErrInternal("pubsub.Message.As returned failed").
					WithMeta("msgId", msg.LoggableID)
			}
			// Get JetStream metadata.
			var err error
			meta, err = natsMsg.Metadata()
			if err != nil {
				return err
			}
		}

		// Set up the timestamp in the context.
		ctx = xpubsub.CtxSetMessageSequence(ctx, int64(meta.Sequence.Stream))

		// Handle the next handlers.
		if err := next.Handle(ctx, msg); err != nil {
			return err
		}
		return nil
	})
}
