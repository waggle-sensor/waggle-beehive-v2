package main

import (
	"context"
)

// DummyBackend implements a test backend which replays a list of test records.
type DummyBackend struct {
	Records []*Record
}

// Query provides the test records through Results.
func (backend *DummyBackend) Query(ctx context.Context, query *Query) (Results, error) {
	return &dummyResults{records: backend.Records}, nil
}

type dummyResults struct {
	records []*Record
	record  *Record
}

func (r *dummyResults) Err() error {
	return nil
}

func (r *dummyResults) Close() error {
	return nil
}

func (r *dummyResults) Record() *Record {
	return r.record
}

func (r *dummyResults) Next() bool {
	if len(r.records) == 0 {
		return false
	}
	r.record = r.records[0]
	r.records = r.records[1:]
	return true
}
