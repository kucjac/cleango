package xpubsub

type subscriptionCtx string

var (
	// SubTopicCtxKey is the key for the subscription topic.
	SubTopicCtxKey = subscriptionCtx("sub:topic")
	// SubIdCtxKey is the key for the subscription id.
	SubIdCtxKey = subscriptionCtx("sub:id")
)
