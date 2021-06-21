package xpubsub

type subscriptionCtx string

var (
	// SubscriptionSubjectCtxKey is the key for the subscription topic.
	SubscriptionSubjectCtxKey = subscriptionCtx("sub:topic")
	// SubscriptionIdCtxKey is the key for the subscription id.
	SubscriptionIdCtxKey = subscriptionCtx("sub:id")
)
