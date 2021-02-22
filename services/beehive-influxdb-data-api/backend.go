package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	influxdb2query "github.com/influxdata/influxdb-client-go/v2/api/query"
)

// Backend defines an interface for query backends.
type Backend interface {
	Query(context.Context, *Query) (Results, error)
}

// Results defines an interface for query result sets.
type Results interface {
	Err() error
	Close() error
	Next() bool
	Record() *Record
}

// Record holds an SDR API record.
type Record struct {
	Timestamp time.Time         `json:"timestamp"`
	Name      string            `json:"name"`
	Value     interface{}       `json:"value"`
	Meta      map[string]string `json:"meta"`
}

// InfluxBackend provides backend to InfluxDB.
type InfluxBackend struct {
	Client influxdb2.Client
	Bucket string
}

// Query converts and makes a query to an Influx backend.
func (backend *InfluxBackend) Query(ctx context.Context, query *Query) (Results, error) {
	fluxQuery, err := buildFluxQuery(backend.Bucket, query)
	if err != nil {
		return nil, err
	}

	queryAPI := backend.Client.QueryAPI(backend.Bucket)

	results, err := queryAPI.Query(ctx, fluxQuery)
	if err != nil {
		return nil, err
	}

	return &InfluxResults{results: results}, nil
}

// InfluxResults manages an InfluxDB result set.
type InfluxResults struct {
	results *api.QueryTableResult
	record  *Record
	err     error
}

// Err returns the last error caught
func (r *InfluxResults) Err() error {
	return r.err
}

// Close closes the record set
func (r *InfluxResults) Close() error {
	return r.results.Close()
}

// Record returns the current record
func (r *InfluxResults) Record() *Record {
	return r.record
}

// Next gets the next record, if available
func (r *InfluxResults) Next() bool {
	if !r.results.Next() {
		return false
	}

	r.record, r.err = convertToAPIRecord(r.results.Record())
	if r.err != nil {
		return false
	}

	return true
}

func convertToAPIRecord(rec *influxdb2query.FluxRecord) (*Record, error) {
	apirec := &Record{}

	name, ok := rec.Values()["_measurement"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid measurement name type")
	}

	apirec.Name = name
	apirec.Timestamp = rec.Time()
	apirec.Value = rec.Values()["_value"]
	apirec.Meta = buildMetaFromRecord(rec)
	return apirec, nil
}

func buildMetaFromRecord(rec *influxdb2query.FluxRecord) map[string]string {
	meta := make(map[string]string)

	for k, v := range rec.Values() {
		// skip influxdb internal fields (convention is to start with _)
		if strings.HasPrefix(k, "_") {
			continue
		}
		// skip influxdb "leaky" fields (table and result don't use above
		// convention but still leak internal details about query)
		if k == "table" || k == "result" {
			continue
		}
		// only include string types in meta
		s, ok := v.(string)
		if !ok {
			continue
		}
		meta[k] = s
	}

	return meta
}
