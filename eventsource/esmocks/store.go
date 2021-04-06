package esmocks

import (
	"github.com/golang/mock/gomock"

	"github.com/kucjac/cleango/eventsource"
	"github.com/kucjac/cleango/eventsource/internal/storemock"
	"github.com/kucjac/cleango/messages/codec"
)

var _ eventsource.EventStore = (*MockStore)(nil)

// MockStore is a mocked eventsource.EventStore with predefined aggregate base setter.
type MockStore struct {
	*storemock.MockStore
	eventsource.AggregateBaseSetter
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller, eventCodec, snapCodec codec.Codec, idGen eventsource.IdGenerator) *MockStore {
	return &MockStore{
		MockStore:           storemock.NewMockStore(ctrl),
		AggregateBaseSetter: eventsource.NewAggregateBaseSetter(eventCodec, snapCodec, idGen),
	}
}

// NewDefaultMockStore creates new default mock store.
func NewDefaultMockStore(ctrl *gomock.Controller) *MockStore {
	c := codec.JSON()
	return &MockStore{
		MockStore:           storemock.NewMockStore(ctrl),
		AggregateBaseSetter: eventsource.NewAggregateBaseSetter(c, c, eventsource.UUIDGenerator{}),
	}
}
