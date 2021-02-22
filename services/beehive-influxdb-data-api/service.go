package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2query "github.com/influxdata/influxdb-client-go/v2/api/query"
)

// Service keeps the service configuration for the SDR API service.
type Service struct {
	Client influxdb2.Client
	Bucket string
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

	query, err := parseAPIQuery(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	fluxQuery, err := buildFluxQuery(svc.Bucket, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	queryStart := time.Now()
	queryCount := 0

	queryAPI := svc.Client.QueryAPI(svc.Bucket)

	// run query against influxdb and get results
	results, err := queryAPI.Query(r.Context(), fluxQuery)
	if err != nil {
		log.Printf("influxdb query error: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: failed to query influxdb"), http.StatusInternalServerError)
		return
	}
	defer results.Close()

	// write ok header before writing results
	w.WriteHeader(http.StatusOK)

	// write all results to client
	for results.Next() {
		rec, err := convertToAPIRecord(results.Record())
		if err != nil {
			log.Printf("invalid influxdb record: %s", err)
			continue
		}
		if err := writeAPIRecord(w, rec); err != nil {
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

func parseAPIQuery(r io.Reader) (*Query, error) {
	var query Query
	if err := json.NewDecoder(r).Decode(&query); err != nil {
		return nil, err
	}
	return &query, nil
}

func writeAPIRecord(w io.Writer, rec *apirecord) error {
	return json.NewEncoder(w).Encode(rec)
}

func convertToAPIRecord(rec *influxdb2query.FluxRecord) (*apirecord, error) {
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
