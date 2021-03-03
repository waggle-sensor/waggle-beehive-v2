package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type serviceError struct {
	Error   error
	Message string
	Code    int
}

// Service keeps the service configuration for the SDR API service.
type Service struct {
	Backend Backend
}

// ServeHTTP dispatches an HTTP request to the right handler.
func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO add simple index page
	switch r.URL.Path {
	case "/api/v1/query":
		if err := svc.serveQuery(w, r); err != nil {
			log.Printf("error %q %v", r.URL.Path, err.Error)
			http.Error(w, err.Message, err.Code)
		} else {
			log.Printf("served %q", r.URL.Path)
		}
	default:
		http.NotFound(w, r)
	}
}

// serveQuery parses a query request, translates and forwards it to InfluxDB
// and writes the results back to the client.
func (svc *Service) serveQuery(w http.ResponseWriter, r *http.Request) *serviceError {
	if r.Method != http.MethodPost {
		return &serviceError{errors.New("invalid method"), "query api expects POST request", http.StatusMethodNotAllowed}
	}

	query, err := parseQuery(r.Body)
	if err != nil {
		return &serviceError{err, err.Error(), http.StatusBadRequest}
	}

	queryStart := time.Now()
	queryCount := 0

	results, err := svc.Backend.Query(r.Context(), query)
	if err != nil {
		return &serviceError{err, "failed to query backend", http.StatusInternalServerError}
	}
	defer results.Close()

	// write all results to client
	w.WriteHeader(http.StatusOK)

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
	return nil
}

func parseQuery(r io.Reader) (*Query, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	query := &Query{}
	if err := decoder.Decode(query); err != nil {
		return nil, err
	}
	return query, nil
}

func writeRecord(w io.Writer, rec *Record) error {
	return json.NewEncoder(w).Encode(rec)
}
