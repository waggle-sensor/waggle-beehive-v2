package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQueryResponse(t *testing.T) {
	records := []*Record{
		{
			Timestamp: time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
			Name:      "sys.uptime",
			Value:     100321,
			Meta: map[string]string{
				"node":   "0000000000000001",
				"plugin": "status:1.0.2",
			},
		},
		{
			Timestamp: time.Date(2022, 1, 1, 10, 30, 0, 0, time.UTC),
			Name:      "env.temp.htu21d",
			Value:     2.3,
			Meta: map[string]string{
				"node":   "0000000000000001",
				"plugin": "metsense:1.0.2",
			},
		},
		{
			Timestamp: time.Date(2023, 2, 1, 10, 45, 0, 0, time.UTC),
			Name:      "raw.htu21d",
			Value:     "234124123",
			Meta: map[string]string{
				"node":   "0000000000000002",
				"plugin": "metsense:1.0.2",
			},
		},
	}

	svc := &Service{
		Backend: &DummyBackend{records},
	}

	body := bytes.NewBufferString(`{
		"start": "-4h"
	}`)

	r := httptest.NewRequest("POST", "/api/v1/query", body)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status ok. got %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)

	// check that output from server is just newline separated json records in same order
	for _, record := range records {
		if !scanner.Scan() {
			t.Fatalf("expected response for record %v", record)
		}
		b1, _ := json.Marshal(record)
		b2 := scanner.Bytes()
		if bytes.Compare(b1, b2) != 0 {
			t.Fatalf("records don't match\nexpect: %s\noutput: %s", b1, b2)
		}
	}
}
