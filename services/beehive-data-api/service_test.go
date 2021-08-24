package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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

	r := httptest.NewRequest("POST", "/", body)
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
		if !bytes.Equal(b1, b2) {
			t.Fatalf("records don't match\nexpect: %s\noutput: %s", b1, b2)
		}
	}
}

func TestQueryDisallowedField(t *testing.T) {
	svc := &Service{
		Backend: &DummyBackend{},
	}

	body := bytes.NewBufferString(`{
		"start": "-4h",
		"filters": {
			"node": "node123"
		}
	}`)

	r := httptest.NewRequest("POST", "/", body)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, r)
	resp := w.Result()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response")
	}

	if !strings.Contains(string(b), `json: unknown field "filters"`) {
		t.Fatalf("expected json error message. got %s.", string(b))
	}

	expectStatus := http.StatusBadRequest
	if resp.StatusCode != expectStatus {
		t.Fatalf("expected status %d. got %d", expectStatus, resp.StatusCode)
	}
}

func TestContentDispositionHeader(t *testing.T) {
	svc := &Service{
		Backend: &DummyBackend{},
	}

	body := bytes.NewBufferString(`{
		"start": "-4h"
	}`)

	r := httptest.NewRequest("POST", "/", body)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, r)
	resp := w.Result()

	pattern := regexp.MustCompile("attachment; filename=\"sage-download-(.+).ndjson\"")

	s := resp.Header.Get("Content-Disposition")

	if !pattern.MatchString(s) {
		t.Fatalf("response must proper Content-Disposition header. got %q", s)
	}
}

func TestNoPayload(t *testing.T) {
	svc := &Service{
		Backend: &DummyBackend{},
	}

	r := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, r)
	resp := w.Result()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "error: must provide a request body\n" {
		t.Fatalf("missing error message when missing request body")
	}
}
