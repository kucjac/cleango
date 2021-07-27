package eventstate

import (
	"testing"
	"time"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/codec"
	"github.com/kucjac/cleango/database/es"
)

const (
	testEventId   = "ca526187-c448-4514-9979-dab2ffb38fd7"
	testEventType = "test_event_type"
	testAggregate = "test_aggregate"
	testHandler1  = "test_handler_1"
	testHandler2  = "test_handler_2"
)

var (
	testEvent = es.Event{
		EventId:       testEventId,
		EventType:     testEventType,
		AggregateType: testAggregate,
		AggregateId:   "3f4c21e2-0b40-49a1-ab89-dcdbbe4630c7",
		Timestamp:     22608,
		Revision:      1,
	}
	bs = es.NewAggregateBaseSetter(codec.JSON(), codec.JSON(), es.UUIDGenerator{})
)

func TestEventState(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		e, err := InitializeUnhandledEventState(&testEvent, bs, nil)
		if err != nil {
			t.Fatalf("initialize event state failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("start handling failed: %v", err)
		}

		err = e.FinishHandling(testHandler1)
		if err != nil {
			t.Fatalf("finish handling failed: %v", err)
		}

		h := e.handlers[testHandler1]
		if h.latestState != StateFinished {
			t.Errorf("latest state of handler is different than FINISHED - %v", h.latestState)
		}
		if h.finishedAt.IsZero() {
			t.Error("event state for handler is not marked as finished")
		}
		if len(h.handles) != 2 {
			t.Error("more handles defined in the handler than expected")
		}
	})

	t.Run("DoubleHandling", func(t *testing.T) {
		e, err := InitializeUnhandledEventState(&testEvent, bs, nil)
		if err != nil {
			t.Fatalf("initialize event state failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("start handling failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err == nil {
			t.Fatalf("event state should fail for second handle in row")
		}
	})

	t.Run("StartAfterFinish", func(t *testing.T) {
		e, err := InitializeUnhandledEventState(&testEvent, bs, nil)
		if err != nil {
			t.Fatalf("initialize event state failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("start handling failed: %v", err)
		}

		err = e.FinishHandling(testHandler1)
		if err != nil {
			t.Fatalf("finish handling failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err == nil {
			t.Fatalf("event state should fail for second handle in row")
		}
	})

	t.Run("MultipleHandlers", func(t *testing.T) {
		e, err := InitializeUnhandledEventState(&testEvent, bs, nil)
		if err != nil {
			t.Fatalf("initialize event state failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("start handling failed: %v", err)
		}

		err = e.FinishHandling(testHandler1)
		if err != nil {
			t.Fatalf("finish handling failed: %v", err)
		}

		err = e.StartHandling(testHandler2)
		if err != nil {
			t.Fatalf("starting handling event for different state failed")
		}

		err = e.FinishHandling(testHandler2)
		if err != nil {
			t.Fatalf("finish handling failed: %v", err)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		o := DefaultOptions()
		o.MinFailInterval = time.Millisecond
		o.MaxFailures = 1

		e, err := InitializeUnhandledEventState(&testEvent, bs, o)
		if err != nil {
			t.Fatalf("initialize event state failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("start handling failed: %v", err)
		}

		err = e.HandlingFailed(testHandler1, cgerrors.ErrInvalidArgument("example error"))
		if err != nil {
			t.Fatalf("failing handling failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err == nil {
			t.Fatal("handling should failed as it occurred too fast")
		}

		time.Sleep(o.MinFailInterval)
		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("new handle should not fail - %v", err)
		}

		err = e.HandlingFailed(testHandler1, cgerrors.ErrInvalidArgument("example error"))
		if err != nil {
			t.Fatalf("failing handling failed: %v", err)
		}

		// Sleep for minimum interval.
		time.Sleep(o.MinFailInterval * 2)
		err = e.StartHandling(testHandler1)
		if err == nil {
			t.Fatal("starting third handling should fail - too many failures")
		}
		err = e.ResetFailures(testHandler1)
		if err != nil {
			t.Fatalf("reseting failures failed: %v", err)
		}

		err = e.StartHandling(testHandler1)
		if err != nil {
			t.Fatalf("starting handling failed after reset: %v", err)
		}
	})
}
