package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var allowedTags = []string{"node", "plugin", "camera"}

// Service keeps the service configuration for the SDR API service.
type Service struct {
	Client influxdb2.Client
	Bucket string
}

func (svc *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var query Query

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	fluxQuery, err := buildFluxQuery(svc.Bucket, &query)
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Printf("query: %s", fluxQuery)

	queryAPI := svc.Client.QueryAPI(svc.Bucket)

	queryStart := time.Now()
	queryCount := 0

	results, err := queryAPI.Query(context.Background(), fluxQuery)
	if err != nil {
		log.Printf("query error: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: failed to query influxdb"), http.StatusInternalServerError)
		return
	}
	defer results.Close()

	// go ahead and write header
	w.WriteHeader(http.StatusOK)

	for results.Next() {
		rec := &struct {
			Timestamp time.Time         `json:"timestamp"`
			Name      string            `json:"name"`
			Value     interface{}       `json:"value"`
			Meta      map[string]string `json:"meta"`
		}{
			Meta: make(map[string]string),
		}

		values := results.Record().Values()

		if s, ok := values["_measurement"].(string); ok {
			rec.Name = s
		} else {
			log.Printf("invalid measurement name")
			continue
		}

		rec.Timestamp = results.Record().Time()
		rec.Value = values["_value"]

		// populate meta fields
		for _, k := range allowedTags {
			v, ok := values[k]
			if !ok {
				continue
			}
			s, ok := v.(string)
			if !ok {
				continue
			}
			rec.Meta[k] = s
		}

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
