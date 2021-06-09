package esxsql_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
	"github.com/kucjac/cleango/database/es/esxsql"
	"github.com/kucjac/cleango/database/xpq"
	"github.com/kucjac/cleango/database/xsql"
	_ "github.com/lib/pq"
)

func testPostgresStore(t *testing.T) *esxsql.Storage {
	conn := testPostgresConn(t)
	config := esxsql.DefaultConfig()
	config.SchemaName = "public"
	s, err := esxsql.New(conn, config, xpq.NewDriver())
	if err != nil {
		t.Fatalf("creating esxsql storage failed: %v", err)
	}
	return s
}

func testTx(t *testing.T, s *esxsql.Storage) (*esxsql.Transaction, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	tx, err := s.BeginTx(ctx)
	if err != nil {
		t.Fatalf("starting transaction failed: %v", err)
	}
	return tx, func() {
		tx.Rollback()
		cancel()
	}
}

func testPostgresConn(t *testing.T) *xsql.Conn {
	uri := os.Getenv("CG_PG_TEST_URI")
	if uri == "" {
		t.Skip("no CG_PG_TEST_URI defined...")
	}
	conn, err := xsql.Connect("postgres", uri)
	if err != nil {
		t.Fatalf("establishing postgres connection failed: %v", err)
	}
	return conn
}

func now() int64 { return time.Now().UTC().UnixNano() }

const (
	aggId     = "12430e53-edba-476a-ac07-989b201b89e4"
	agg2ID    = "f2b19e0f-cb6c-4dc3-8c3a-341762d4b87f"
	eventType = "TEST_EVENT_TYPE"
	aggType   = "TEST_AGG_TYPE"
)

var (
	e1 = es.Event{
		EventId:       "0a76941b-08ec-4bb9-bae5-7b8d8f6623b6",
		EventType:     eventType,
		AggregateType: aggType,
		AggregateId:   aggId,
		EventData:     []byte(`{"name":"some name"}`),
		Timestamp:     now(),
		Revision:      1,
	}
	e2 = es.Event{
		EventId:       "45606771-e303-496b-8399-f1a71762735c",
		EventType:     "OTHER_EVENT_TYPE",
		AggregateType: aggType,
		AggregateId:   aggId,
		EventData:     nil,
		Timestamp:     now(),
		Revision:      2,
	}
	e3 = es.Event{
		EventId:       "4cedbacb-3480-4499-b977-f6b0aaaa5ad1",
		EventType:     "EVENT_TYPE",
		AggregateType: aggType,
		AggregateId:   agg2ID,
		EventData:     []byte(`{"name":"some name"}`),
		Timestamp:     now(),
		Revision:      1,
	}
	e4 = es.Event{
		EventId:       "472766cf-02ce-4a23-809a-01686718bec6",
		EventType:     "EVENT_TYPE",
		AggregateType: aggType,
		AggregateId:   "c37aea31-c602-46c7-8ded-b172a4b81ace",
		EventData:     []byte(`{"name":"some name"}`),
		Timestamp:     now(),
		Revision:      1,
	}
	e5 = es.Event{
		EventId:       "951d6648-371d-4d57-a065-202afd985d29",
		EventType:     "EVENT_TYPE",
		AggregateType: aggType,
		AggregateId:   aggId,
		EventData:     nil,
		Timestamp:     now(),
		Revision:      3,
	}
)

func TestPostgresEvents(t *testing.T) {
	store := testPostgresStore(t)

	ctx := context.Background()
	t.Run("Batch", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()

		err := tx.SaveEvents(ctx, []*es.Event{&e1, &e2, &e3, &e4, &e5})
		if err != nil {
			t.Fatalf("saving events failed: %v", err)
		}

		events, err := tx.ListEvents(ctx, aggId, aggType)
		if err != nil {
			t.Fatalf("listing events failed: %v", err)
		}

		if len(events) != 3 {
			t.Fatalf("result should contain three events, has: %d", len(events))
		}

		for i, e := range events {
			var expectedEvent *es.Event
			switch i {
			case 0:
				expectedEvent = &e1
			case 1:
				expectedEvent = &e2
			case 2:
				expectedEvent = &e5
			}
			compareEvents(t, e, expectedEvent, i)
		}

		events, err = tx.ListEventsAfterRevision(ctx, aggId, aggType, e1.Revision)
		if err != nil {
			t.Fatalf("listing events from revision failed: %v", err)
		}

		if len(events) != 2 {
			t.Fatalf("there should be exactly two events starting from revision '2' but there are: %d", len(events))
		}

		for i, e := range events {
			var expectedEvent *es.Event
			switch i {
			case 0:
				expectedEvent = &e2
			case 1:
				expectedEvent = &e5
			}
			compareEvents(t, e, expectedEvent, i)
		}
	})

	t.Run("Single", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()

		err := tx.SaveEvents(ctx, []*es.Event{&e3})
		if err != nil {
			t.Fatalf("saving single event failed: %v", err)
		}

		events, err := tx.ListEvents(ctx, e3.AggregateId, e3.AggregateType)
		if err != nil {
			t.Fatalf("getting single event failed: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("")
		}

		e := events[0]
		compareEvents(t, e, &e3, -1)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()

		err := tx.SaveEvents(ctx, []*es.Event{&e4})
		if err != nil {
			t.Fatalf("saving single event failed: %v", err)
		}
		err = tx.SaveEvents(ctx, []*es.Event{&e4})
		if err == nil {
			t.Error("saving single duplicated event should fail")
		} else {
			if tx.ErrorCode(err) != cgerrors.ErrorCode_AlreadyExists {
				t.Errorf("saving single duplicated event should return error of type AlreadyExists but is: %v", err)
			}
		}
	})

	t.Run("Stream", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()

		err := tx.SaveEvents(ctx, []*es.Event{&e1, &e2, &e3, &e4, &e5})
		if err != nil {
			t.Fatalf("saving events failed: %v", err)
		}

		stream, err := tx.StreamEvents(ctx, &es.StreamEventsRequest{BuffSize: 2})
		if err != nil {
			t.Fatalf("getting event stream failed: %v", err)
		}

		var i int
		for e := range stream {
			var expected *es.Event
			switch i {
			case 0:
				expected = &e1
			case 1:
				expected = &e2
			case 2:
				expected = &e3
			case 3:
				expected = &e4
			case 4:
				expected = &e5
			}
			compareEvents(t, e, expected, i)
			i++
		}

		if i != 5 {
			t.Errorf("obtained fewer events: %d", i)
		}
	})
}

func TestPostgresSnapshots(t *testing.T) {
	store := testPostgresStore(t)
	ctx := context.Background()
	t.Run("AlreadyExists", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()
		snap := &es.Snapshot{
			AggregateId:      aggId,
			AggregateType:    aggType,
			AggregateVersion: 1,
			Revision:         1,
			Timestamp:        now(),
			SnapshotData:     nil,
		}

		err := tx.SaveSnapshot(ctx, snap)
		if err != nil {
			t.Fatalf("saving snapshot failed: %v", err)
		}

		err = tx.SaveSnapshot(ctx, snap)
		if err == nil {
			t.Error("expected error already exists on saving duplicated snapshot")
		} else {
			if code := tx.ErrorCode(err); code != cgerrors.ErrorCode_AlreadyExists {
				t.Errorf("expected error already exists on saving duplicated snapshot, but got: %v, %v", code.String(), err)
			}
		}
	})

	t.Run("Valid", func(t *testing.T) {
		tx, cf := testTx(t, store)
		defer cf()
		snap := &es.Snapshot{
			AggregateId:      aggId,
			AggregateType:    aggType,
			AggregateVersion: 1,
			Revision:         1,
			Timestamp:        now(),
			SnapshotData:     nil,
		}

		err := tx.SaveSnapshot(ctx, snap)
		if err != nil {
			t.Fatalf("saving snapshot failed: %v", err)
		}

		taken, err := tx.GetSnapshot(ctx, aggId, aggType, 1)
		if err != nil {
			t.Fatalf("getting snapshot failed: %v", err)
		}

		compareSnapshots(t, taken, snap)
	})
}

func compareEvents(t *testing.T, e *es.Event, expectedEvent *es.Event, i int) {
	if e.EventId != expectedEvent.EventId {
		t.Errorf("event at index: %d mismatch value of EventId, is: %v, want: %v", i, e.EventId, expectedEvent.EventId)
	}
	if e.EventType != expectedEvent.EventType {
		t.Errorf("event at index: %d mismatch value of EventType, is: %v, want: %v", i, e.EventType, expectedEvent.EventType)
	}
	if e.AggregateType != expectedEvent.AggregateType {
		t.Errorf("event at index: %d mismatch value of AggregateType, is: %v, want: %v", i, e.AggregateType, expectedEvent.AggregateType)
	}
	if e.AggregateId != expectedEvent.AggregateId {
		t.Errorf("event at index: %d mismatch value of AggregateId, is: %v, want: %v", i, e.AggregateId, expectedEvent.AggregateId)
	}
	if !bytes.Equal(e.EventData, expectedEvent.EventData) {
		t.Errorf("event at index: %d mismatch value of EventData, is: %v, want: %v", i, e.EventData, expectedEvent.EventData)
	}
	if e.Timestamp != expectedEvent.Timestamp {
		t.Errorf("event at index: %d mismatch value of Timestamp, is: %v, want: %v", i, e.Timestamp, expectedEvent.Timestamp)
	}
	if e.Revision != expectedEvent.Revision {
		t.Errorf("event at index: %d mismatch value of Revision, is: %v, want: %v", i, e.Revision, expectedEvent.Revision)
	}
}

func compareSnapshots(t *testing.T, s, compare *es.Snapshot) {
	if s.AggregateId != compare.AggregateId {
		t.Errorf("snapshot AggregateId: %v different than expected: %v", s.AggregateId, compare.AggregateId)
	}
	if s.AggregateType != compare.AggregateType {
		t.Errorf("snapshot AggregateType: %v different than expected: %v", s.AggregateType, compare.AggregateType)
	}
	if s.AggregateVersion != compare.AggregateVersion {
		t.Errorf("snapshot AggregateVersion: %v different than expected: %v", s.AggregateVersion, compare.AggregateVersion)
	}
	if s.Timestamp != compare.Timestamp {
		t.Errorf("snapshot Timestamp: %v different than expected: %v", s.Timestamp, compare.Timestamp)
	}
	if !bytes.Equal(s.SnapshotData, compare.SnapshotData) {
		t.Errorf("snapshot SnapshotData: %x different than expected: %x", s.SnapshotData, compare.SnapshotData)
	}
	if s.Revision != compare.Revision {
		t.Errorf("snapshot Revision: %v different than expected: %v", s.Revision, compare.Revision)
	}
}
