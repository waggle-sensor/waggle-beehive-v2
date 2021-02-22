package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2query "github.com/influxdata/influxdb-client-go/v2/api/query"
)

var allowedTags = []string{"node", "plugin", "camera"}

// Service keeps the service configuration for the SDR API service.
type Service struct {
	Client influxdb2.Client
	Bucket string
}

func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO add simple index page
	switch r.URL.Path {
	case "/api/v1/query":
		svc.serveQuery(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (svc *Service) serveQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "query api expects http POST", http.StatusMethodNotAllowed)
		return
	}

	// parse query body from client
	var query Query

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// build flux query for influxdb
	fluxQuery, err := buildFluxQuery(svc.Bucket, &query)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Printf("query: %s", fluxQuery)

	queryAPI := svc.Client.QueryAPI(svc.Bucket)

	queryStart := time.Now()
	queryCount := 0

	// run query against influxdb and get results
	results, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Printf("query error: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: failed to query influxdb"), http.StatusInternalServerError)
		return
	}
	defer results.Close()

	// write ok header before writing results
	w.WriteHeader(http.StatusOK)

	// write all results to client
	for results.Next() {
		// build api record from influxdb record
		rec, err := buildAPIRecord(results.Record())
		if err != nil {
			log.Printf("invalid influxdb record: %s", err)
			continue
		}

		// write api record to client
		if err := json.NewEncoder(w).Encode(rec); err != nil {
			break
		}

		queryCount++
	}

	if err := results.Err(); err != nil {
		log.Printf("error: %s", err)
	}

	queryDuration := time.Since(queryStart)
	log.Printf("served %d records in %s", queryCount, queryDuration)
}

type apirecord struct {
	Timestamp time.Time         `json:"timestamp"`
	Name      string            `json:"name"`
	Value     interface{}       `json:"value"`
	Meta      map[string]string `json:"meta"`
}

func buildAPIRecord(rec *influxdb2query.FluxRecord) (*apirecord, error) {
	apirec := &apirecord{}

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

	for _, k := range allowedTags {
		v, ok := rec.Values()[k]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		meta[k] = s
	}

	return meta
}
