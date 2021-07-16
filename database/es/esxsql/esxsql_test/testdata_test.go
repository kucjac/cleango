package esxsql_tst

import (
	"github.com/kucjac/cleango/database/es"
)

const (
	aggId          = "12430e53-edba-476a-ac07-989b201b89e4"
	agg2ID         = "f2b19e0f-cb6c-4dc3-8c3a-341762d4b87f"
	eventType      = "TEST_EVENT_TYPE"
	otherEventType = "OTHER_EVENT_TYPE"
	aggType        = "TEST_AGG_TYPE"
	testHandler    = "TEST_HANDLER"
	testHandler2   = "TEST_HANDLER_2"
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
