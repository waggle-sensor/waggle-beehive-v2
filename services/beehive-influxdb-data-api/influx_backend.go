package main

import (
	"context"
	"fmt"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	influxdb2query "github.com/influxdata/influxdb-client-go/v2/api/query"
)

// InfluxBackend implements a backend to InfluxDB.
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

	return &influxResults{results: results}, nil
}

type influxResults struct {
	results *api.QueryTableResult
	record  *Record
	err     error
}

func (r *influxResults) Err() error {
	return r.err
}

func (r *influxResults) Close() error {
	return r.results.Close()
}

func (r *influxResults) Record() *Record {
	return r.record
}

func (r *influxResults) Next() bool {
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
