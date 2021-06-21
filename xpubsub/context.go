package xpubsub

import (
	"context"
	"time"
)

type subscriptionCtx string

var (
	// SubscriptionSubjectCtxKey is the key for the subscription topic.
	SubscriptionSubjectCtxKey = subscriptionCtx("pubsub:subject")
	// SubscriptionIdCtxKey is the key for the subscription id.
	SubscriptionIdCtxKey = subscriptionCtx("pubsub:id")
	// MessageTimestampKey is the key for the message timestamp.
	MessageTimestampKey = subscriptionCtx("pubsub:msg.timestamp")
	// MessageSequenceKey is the key for the message timestamp.
	MessageSequenceKey = subscriptionCtx("pubsub:msg.sequence")
)

// CtxSetMessageSequence sets the message sequence in the context,
func CtxSetMessageSequence(ctx context.Context, seq int64) context.Context {
	return context.WithValue(ctx, MessageSequenceKey, seq)
}

// CtxGetMessageSequence gets the message sequence from the context.
func CtxGetMessageSequence(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(MessageSequenceKey).(int64)
	return v, ok
}

// CtxSetMessageTimestamp sets the message timestamp in the context.
func CtxSetMessageTimestamp(ctx context.Context, timestamp time.Time) context.Context {
	return context.WithValue(ctx, MessageTimestampKey, timestamp)
}

// CtxGetMessageTimestamp gets the message timestamp from the context.
func CtxGetMessageTimestamp(ctx context.Context) (time.Time, bool) {
	v, ok := ctx.Value(MessageTimestampKey).(time.Time)
	return v, ok
}
