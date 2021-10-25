package main

import (
	"context"
	"time"
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

// Query holds an SDR query body.
type Query struct {
	Bucket *string           `json:"bucket,omitempty"`
	Start  string            `json:"start,omitempty"`
	End    string            `json:"end,omitempty"`
	Tail   *int              `json:"tail,omitempty"`
	Filter map[string]string `json:"filter"`
}

// Record holds an SDR API record.
type Record struct {
	Timestamp time.Time         `json:"timestamp"`
	Name      string            `json:"name"`
	Value     interface{}       `json:"value"`
	Meta      map[string]string `json:"meta"`
}
