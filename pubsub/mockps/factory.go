package mockps

import (
	"github.com/golang/mock/gomock"
	pubsub2 "github.com/kucjac/cleango/pubsub"

	"github.com/kucjac/cleango/pubsub/codec"
)

// NewMockPublisherFactory creates new mock publisher factory.
func NewMockPublisherFactory(m *gomock.Controller) *MockPublisherFactory {
	return &MockPublisherFactory{m: m}
}

// MockPublisherFactory is the mocked factory for the messages publisher.
type MockPublisherFactory struct {
	m *gomock.Controller
}

// NewPublisher creates new mocked publisher.
func (m MockPublisherFactory) NewPublisher(codec.Codec) (pubsub2.Publisher, error) {
	return NewMockPublisher(m.m), nil
}

// NewMockSubscriberFactory creates new mock subscriber factory.
func NewMockSubscriberFactory(m *gomock.Controller) *MockSubscriberFactory {
	return &MockSubscriberFactory{m: m}
}

// MockSubscriberFactory is the mocked factory for the messages subscriber.
type MockSubscriberFactory struct {
	m *gomock.Controller
}

// NewSubscriber creates new mmocked subscriber.
func (m *MockSubscriberFactory) NewSubscriber(_ ...pubsub2.SubscriptionOption) (pubsub2.Subscriber, error) {
	return NewMockSubscriber(m.m), nil
}
