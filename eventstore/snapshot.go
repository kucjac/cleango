package eventstore

// Snapshot is a structure that define basic fields of the aggregate snapshot.
type Snapshot struct {
	AggregateId      string `json:"aggregate_id,omitempty"`
	AggregateType    string `json:"aggregate_type,omitempty"`
	AggregateVersion int64  `json:"aggregate_version,omitempty"`
	Revision         int64  `json:"revision,omitempty"`
	Timestamp        int64  `json:"timestamp,omitempty"`
	SnapshotData     []byte `json:"snapshot_data,omitempty"`
}
