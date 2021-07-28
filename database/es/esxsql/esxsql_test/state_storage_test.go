package esxsql_tst

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database/es"
	"github.com/kucjac/cleango/database/es/esstate"
	"github.com/kucjac/cleango/database/es/esxsql"
	"github.com/kucjac/cleango/database/xsql"
	"github.com/kucjac/cleango/ddd/events/eventstate"
)

func testPostgresStateStore(t *testing.T) *esxsql.StateStorage {
	schemaName := esxsql.ToSnakeCase(t.Name())
	schemaName = strings.ReplaceAll(schemaName, "/", "_")
	conn := testPostgresConn(t)
	config := esxsql.DefaultConfig()
	config.SchemaName = schemaName
	config.AggregateTypes = []string{aggType}
	config.WithEventState(esxsql.DefaultEventStateConfig(
		eventstate.Handler{Name: testHandler, EventTypes: []string{eventType, otherEventType}},
		eventstate.Handler{Name: testHandler2, EventTypes: []string{eventType}},
	))
	s, err := esxsql.NewStateStorage(conn, config)
	if err != nil {
		t.Fatalf("creating esxsql storage failed: %v", err)
	}
	return s
}

func testStateTx(t *testing.T, s *esxsql.StateStorage) (esstate.TxStorage, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	tx, err := s.BeginTx(ctx)
	if err != nil {
		t.Fatalf("starting transaction failed: %v", err)
	}

	var txc *xsql.Tx
	err = tx.As(&txc)
	if err != nil {
		t.Fatalf("getting tx conn failed: %v", err)
	}
	config := s.Config()
	if _, err = txc.Exec(fmt.Sprintf("CREATE SCHEMA %s;", config.SchemaName)); err != nil {
		t.Fatalf("creating schema failed: %v", err)
	}

	if err = esxsql.Migrate(txc, &config); err != nil {
		t.Fatalf("migrating failed: %v", err)
	}
	return tx, func() {
		tx.Rollback(context.Background())
		cancel()
	}
}

func TestPostgresStateStorage(t *testing.T) {
	ctx := context.Background()

	store := testPostgresStateStore(t)
	tx, cf := testStateTx(t, store)
	defer cf()

	err := tx.MarkUnhandled(ctx, e1.EventId, e1.EventType, e1.Timestamp)
	if err != nil {
		t.Fatalf("marking unhandled failed: %v", err)
	}

	err = tx.MarkUnhandled(ctx, e2.EventId, e2.EventType, e2.Timestamp)
	if err != nil {
		t.Fatalf("marking unhandled failed: %v", err)
	}

	err = tx.SaveEvents(ctx, []*es.Event{&e1, &e2})
	if err != nil {
		t.Fatalf("saving events failed: %v", err)
	}

	events, err := tx.FindUnhandled(ctx, eventstate.FindUnhandledQuery{HandlerNames: []string{testHandler}})
	if err != nil {
		t.Fatalf("finding unhandled events failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("there should be exactly 2 unhandled events but is: %d", len(events))
	}

	events, err = tx.FindUnhandled(ctx, eventstate.FindUnhandledQuery{HandlerNames: []string{testHandler2}})
	if err != nil {
		t.Fatalf("finding unhandled events failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("there should be exactly 1 unhandled events but is: %d", len(events))
	}

	events, err = tx.FindUnhandled(ctx, eventstate.FindUnhandledQuery{HandlerNames: []string{testHandler, testHandler2}})
	if err != nil {
		t.Fatalf("finding unhandled events failed: %v", err)
	}

	if len(events) != 3 {
		t.Fatalf("there should be exactly 3 unhandled events but is: %d", len(events))
	}

	err = tx.StartHandling(ctx, e1.EventId, testHandler, now())
	if err != nil {
		t.Fatalf("starting handling failed: %v", err)
	}

	events, err = tx.FindUnhandled(ctx, eventstate.FindUnhandledQuery{HandlerNames: []string{testHandler}})
	if err != nil {
		t.Fatalf("finding unhandled failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("there should be only one unhandled event but is: %d", len(events))
	}

	err = tx.HandlingFailed(ctx, &eventstate.HandleFailure{
		EventID:     e1.EventId,
		HandlerName: testHandler,
		Err:         "failed",
		ErrCode:     cgerrors.ErrorCode_Internal,
		RetryNo:     1,
		Timestamp:   time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("handling failure failed: %v", err)
	}

	failed, err := tx.FindFailures(ctx, eventstate.FindFailureQuery{HandlerNames: []string{testHandler, testHandler2}})
	if err != nil {
		t.Fatalf("finding failures failed: %v", err)
	}

	if len(failed) != 1 {
		t.Fatalf("there should be only one unhandled event but is: %d", len(failed))
	}

	err = tx.FinishHandling(ctx, e1.EventId, testHandler, now())
	if err != nil {
		t.Fatalf("finishing handling failed: %v", err)
	}

}
