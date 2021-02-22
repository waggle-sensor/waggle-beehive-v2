package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Service keeps the service configuration for the SDR API service.
type Service struct {
	Backend Backend
}

// apirecord is used to hold a response record in the SDR API format.
type apirecord struct {
	Timestamp time.Time         `json:"timestamp"`
	Name      string            `json:"name"`
	Value     interface{}       `json:"value"`
	Meta      map[string]string `json:"meta"`
}

// ServeHTTP dispatches an HTTP request to the right handler.
func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO add simple index page
	switch r.URL.Path {
	case "/api/v1/query":
		svc.serveQuery(w, r)
	default:
		http.NotFound(w, r)
	}
}

// serveQuery parses a query request, translates and forwards it to InfluxDB
// and writes the results back to the client.
func (svc *Service) serveQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "query api expects http POST", http.StatusMethodNotAllowed)
		return
	}

	query, err := parseQuery(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	queryStart := time.Now()
	queryCount := 0

	results, err := svc.Backend.Query(r.Context(), query)
	if err != nil {
		log.Printf("influxdb query error: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: failed to query influxdb"), http.StatusInternalServerError)
		return
	}
	defer results.Close()

	// write all results to client
	for results.Next() {
		if err := writeRecord(w, results.Record()); err != nil {
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

func parseQuery(r io.Reader) (*Query, error) {
	var query Query
	if err := json.NewDecoder(r).Decode(&query); err != nil {
		return nil, err
	}
	return &query, nil
}

func writeRecord(w io.Writer, rec *Record) error {
	return json.NewEncoder(w).Encode(rec)
}
