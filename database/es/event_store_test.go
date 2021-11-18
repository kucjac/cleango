package es_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/database/es"
	mockes "github.com/kucjac/cleango/database/es/mock"
)

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := es.DefaultConfig()

	storage := mockes.NewMockStorage(ctrl)

	store, err := es.New(cfg, codec.JSON(), codec.JSON(), storage)
	if err != nil {
		t.Fatalf("creating new event storage failed: %v", err)
	}

	ctx := context.Background()

	const aggId = "ad84c877-7e5e-4bb8-a1ca-abb02c48fd0a"
	e1 := &es.Event{
		EventId:       "341abd56-cfbe-4033-9dec-b45ca8cf6c2d",
		EventType:     aggregateCreatedType,
		AggregateType: aggregateType,
		AggregateId:   aggId,
		EventData:     []byte("{}"),
		Timestamp:     now(),
		Revision:      1,
	}

	jsonCreatedAt, err := e1.Time().MarshalJSON()
	if err != nil {
		t.Fatalf("marshaling time to json failed: %v", err)
	}

	e2 := &es.Event{
		EventId:       "5c0bd35c-d090-46a4-bf0d-3b9ebcd2a1c7",
		EventType:     aggregateNameChangedType,
		AggregateType: aggregateType,
		AggregateId:   aggId,
		EventData:     []byte(`{"name": "NewName"}`),
		Timestamp:     now(),
		Revision:      2,
	}

	t.Run("LoadEvents", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			agg := &testAggregate{}
			store.SetAggregateBase(agg, aggId, aggregateType, 1)

			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{e1, e2}, nil)

			if err = store.LoadEvents(ctx, agg); err != nil {
				t.Fatalf("loading events failed: %v", err)
			}

			validateTestAggregate(t, agg, e1, "NewName")
		})

		t.Run("NotFound", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{}, nil)

			err = store.LoadEvents(ctx, agg)
			if err == nil {
				t.Errorf("expected error NotFound, but is nil")
			} else if !cgerrors.IsNotFound(err) {
				t.Errorf("expected error with code NotFound but is: %v", err)
			}
		})
	})

	t.Run("LoadEventsWithSnapshot", func(t *testing.T) {
		t.Run("NoSnapshot", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			// The snapshot error is not found.
			storage.EXPECT().
				ErrorCode(gomock.Any()).
				Return(cgerrors.CodeNotFound)

			storage.EXPECT().
				GetSnapshot(ctx, aggId, aggregateType, int64(1)).
				Return(nil, errors.New("not found"))

			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{e1, e2}, nil)

			if err = store.LoadEventsWithSnapshot(ctx, agg); err != nil {
				t.Errorf("expected no error but is: %v", err)
			}

			validateTestAggregate(t, agg, e1, "NewName")
		})

		t.Run("NotFound", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			// The snapshot error is not found.
			storage.EXPECT().
				ErrorCode(gomock.Any()).
				Return(cgerrors.CodeNotFound)

			storage.EXPECT().
				GetSnapshot(ctx, aggId, aggregateType, int64(1)).
				Return(nil, errors.New("not found"))

			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{}, nil)

			if err = store.LoadEventsWithSnapshot(ctx, agg); err == nil {
				t.Errorf("expected no error but is: %v", err)
			} else if !cgerrors.IsNotFound(err) {
				t.Errorf("expected LoadEventsWithSnapshot errot to be not found but is: %v", err)
			}
		})

		t.Run("SnapshotNoEvents", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			storage.EXPECT().
				GetSnapshot(ctx, aggId, aggregateType, int64(1)).
				Return(&es.Snapshot{
					AggregateId:      aggId,
					AggregateType:    aggregateType,
					AggregateVersion: 1,
					Revision:         2,
					Timestamp:        e2.Timestamp,
					SnapshotData:     []byte(fmt.Sprintf(`{"name":"NewName", "createdAt": %s}`, string(jsonCreatedAt))),
				}, nil)

			storage.EXPECT().
				ListEventsAfterRevision(ctx, aggId, aggregateType, int64(2)).
				Return([]*es.Event{}, nil)

			if err = store.LoadEventsWithSnapshot(ctx, agg); err != nil {
				t.Fatalf("expected no error but is: %v", err)
			}

			validateTestAggregate(t, agg, e1, "NewName")
		})

		t.Run("Valid", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			storage.EXPECT().
				GetSnapshot(ctx, aggId, aggregateType, int64(1)).
				Return(&es.Snapshot{
					AggregateId:      aggId,
					AggregateType:    aggregateType,
					AggregateVersion: 1,
					Revision:         2,
					Timestamp:        e2.Timestamp,
					SnapshotData:     []byte(fmt.Sprintf(`{"name":"NewName", "createdAt": %s}`, string(jsonCreatedAt))),
				}, nil)

			storage.EXPECT().
				ListEventsAfterRevision(ctx, aggId, aggregateType, int64(2)).
				Return([]*es.Event{
					{
						EventId:       "c48da543-d514-4760-981b-bd3e397f4687",
						EventType:     aggregateNameChangedType,
						AggregateType: aggregateType,
						AggregateId:   aggId,
						EventData:     []byte(`{"name":"ChangedName"}`),
						Timestamp:     now(),
						Revision:      3,
					},
				}, nil)

			if err = store.LoadEventsWithSnapshot(ctx, agg); err != nil {
				t.Fatalf("expected no error but is: %v", err)
			}

			validateTestAggregate(t, agg, e1, "ChangedName")
		})
	})

	t.Run("SaveSnapshot", func(t *testing.T) {
		agg := getTestAggregate(store, aggId)

		storage.EXPECT().
			ListEvents(gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]*es.Event{e1, e2}, nil)

		if err = store.LoadEvents(ctx, agg); err != nil {
			t.Fatalf("loading events failed: %v", err)
		}

		data, err := json.Marshal(agg)
		if err != nil {
			t.Fatalf("marshaling aggregate failed: %v", err)
		}
		storage.EXPECT().
			SaveSnapshot(ctx, &es.Snapshot{
				AggregateId:      aggId,
				AggregateType:    aggregateType,
				AggregateVersion: 1,
				Revision:         2,
				Timestamp:        e2.Timestamp,
				SnapshotData:     data,
			}).Return(nil)

		if err = store.SaveSnapshot(ctx, agg); err != nil {
			t.Errorf("expected no error, but have: %v", err)
		}
	})

	t.Run("Commit", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			if err = agg.Base.SetEvent(&aggregateCreated{}); err != nil {
				t.Fatalf("setting aggregate created message failed: %v", err)
			}
			if err = agg.Base.SetEvent(&aggregateNameChanged{Name: "NewName"}); err != nil {
				t.Fatalf("setting name changed message failed: %v", err)
			}

			base := agg.Base
			if len(base.UncommittedEvents()) != 2 {
				t.Fatalf("aggregate base should contain two uncommitted events but have: %d", len(base.UncommittedEvents()))
			}

			storage.EXPECT().
				SaveEvents(ctx, base.UncommittedEvents()).
				Return(nil)

			if err = store.Commit(ctx, agg); err != nil {
				t.Fatalf("committing aggregate messages failed: %v", err)
			}
			base = agg.Base

			// The aggregate should not contain any uncommitted event messages.
			if len(base.UncommittedEvents()) != 0 {
				t.Errorf("an aggregate should not contain any uncommitted events but have: %d", len(base.UncommittedEvents()))
			}

			if len(base.CommittedEvents()) != 2 {
				t.Fatalf("an aggregate should contain two committed events, but have: %d", len(base.CommittedEvents()))
			}
			latestCommitted, ok := base.LatestCommittedEvent()
			if !ok {
				t.Fatalf("latest committed event not found")
			}

			if base.Revision() != 2 {
				t.Errorf("an aggregate should be of revision: '2' but is: '%d'", base.Revision())
			}
			if base.Timestamp() != latestCommitted.Timestamp {
				t.Errorf("an aggregate timestamp should be equal to the last committed event timestamp: '%d', but is: '%d'", latestCommitted.Timestamp, base.Timestamp())
			}
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			agg := getTestAggregate(store, aggId)

			// Assuming that the aggregate was already registered.
			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{e1}, nil)

			if err = store.LoadEvents(ctx, agg); err != nil {
				t.Fatalf("loading events failed: %v", err)
			}

			// Change its name to second name.
			const newName = "SecondName"
			if err = agg.Base.SetEvent(&aggregateNameChanged{Name: newName}); err != nil {
				t.Fatalf("setting name changed message failed: %v", err)
			}

			// Having a situation that another event with revision two for given aggregate.
			storage.EXPECT().
				SaveEvents(ctx, agg.Base.UncommittedEvents()).
				Return(errors.New("event with given revision already exists"))

			storage.EXPECT().ErrorCode(gomock.Any()).Return(cgerrors.CodeAlreadyExists)

			// The snapshot is not found.
			storage.EXPECT().ErrorCode(gomock.Any()).Return(cgerrors.CodeNotFound)
			storage.EXPECT().
				GetSnapshot(ctx, aggId, aggregateType, int64(1)).
				Return(nil, errors.New("snapshot not found"))

			// But it looks like there was already some event with revision 2 for given aggId.
			storage.EXPECT().
				ListEvents(ctx, aggId, aggregateType).
				Return([]*es.Event{e2}, nil)

			// Now it is assumed that it passed to push the change.
			storage.EXPECT().
				SaveEvents(ctx, gomock.Any()).
				Return(nil)

			if err = store.Commit(ctx, agg); err != nil {
				t.Fatalf("committing message should not fail: %v", err)
			}

			if agg.Name != newName {
				t.Errorf("an aggregate is expected to have a name: '%s' but have: '%s", newName, agg.Name)
			}
			base := agg.Base
			if base.Revision() != 3 {
				t.Errorf("an aggregate revision is expected to be: '3' but is: '%d'", base.Revision())
			}

			latest, ok := base.LatestCommittedEvent()
			if !ok {
				t.Fatalf("getting latest committed event failed: %v", err)
			}
			if latest.Revision != 3 {
				t.Errorf("latest committed event should have increased its revision to '3' but have: '%d'", latest.Revision)
			}
		})
	})

	t.Run("StreamEvents", func(t *testing.T) {
		req := &es.StreamEventsRequest{AggregateIDs: []string{aggId}, AggregateTypes: []string{aggregateType}}
		ch := make(chan *es.Event, 2)

		ch <- e1
		ch <- e2
		close(ch)

		storage.EXPECT().
			StreamEvents(ctx, req).
			Return(ch, nil)

		stream, err := store.StreamEvents(ctx, req)
		if err != nil {
			t.Fatalf("streaming events failed: %v", err)
		}

		var i int
		for e := range stream {
			var expected *es.Event
			switch i {
			case 0:
				expected = e1
			case 1:
				expected = e2
			}
			if e != expected {
				t.Errorf("provided different expected event at index: %d", i)
			}
			i++
		}
	})
}

func validateTestAggregate(t *testing.T, agg *testAggregate, e1 *es.Event, name string) {
	if agg.CreatedAt.IsZero() {
		t.Error("aggregate created at should not be zero")
	} else if agg.CreatedAt.UnixNano() != e1.Timestamp {
		t.Errorf("aggregate created at should equal the timestamp in the first event. Expected: %d is : %d", e1.Timestamp, agg.CreatedAt.UnixNano())
	}

	if agg.Name != name {
		t.Errorf("aggregate name should equal to: 'NewName', but is: %s", agg.Name)
	}
}

func getTestAggregate(store *es.Store, aggId string) *testAggregate {
	agg := &testAggregate{}
	store.SetAggregateBase(agg, aggId, aggregateType, 1)
	return agg
}

func now() int64 {
	return time.Now().UTC().UnixNano()
}
