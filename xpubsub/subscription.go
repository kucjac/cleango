package xpubsub

import (
	"gocloud.dev/pubsub"
)

// SubjectSubscription is the structure that contains both the subject and related subscription.
// It is used by the mux to define a human-readable version of the subscription.
// The subject should be the URL formatted subject of specific pubsub driver that matches given subscription.
type SubjectSubscription struct {
	Subject      string
	Subscription *pubsub.Subscription
}
