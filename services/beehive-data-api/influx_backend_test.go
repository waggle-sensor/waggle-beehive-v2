package main

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func init() {
	log.SetOutput(io.Discard)
}

func strptr(s string) *string {
	return &s
}

func intptr(x int) *int {
	return &x
}

func TestBuildFluxQuery(t *testing.T) {
	testcases := []struct {
		Query      *Query
		Expect     string
		ShouldFail bool
	}{
		// this case checks for start in time range
		{
			Query: &Query{
				Start: "-4h",
			},
			Expect: `from(bucket:"mybucket") |> range(start:-4h)`,
		},
		// this case checks for bucket
		{
			Query: &Query{
				Bucket: strptr("downsampled"),
				Start:  "-4h",
				Tail:   intptr(3),
			},
			Expect: `from(bucket:"downsampled") |> range(start:-4h) |> tail(n:3)`,
		},
		// check invalid
		{
			Query: &Query{
				Bucket: strptr("_badbucket"),
				Start:  "-4h",
				Tail:   intptr(3),
			},
			Expect:     ``,
			ShouldFail: true,
		},
		// this case checks for start and end in time range
		{
			Query: &Query{
				Start: "-4h",
				End:   "-2h",
				Tail:  intptr(3),
			},
			Expect: `from(bucket:"mybucket") |> range(start:-4h,stop:-2h) |> tail(n:3)`,
		},
		// this case checks for exact match filter
		{
			Query: &Query{
				Start: "-4h",
				End:   "-2h",
				Filter: map[string]string{
					"node": "0000000000000001",
				}},
			Expect: `from(bucket:"mybucket") |> range(start:-4h,stop:-2h) |> filter(fn: (r) => r.node == "0000000000000001")`,
		},
		// this case checks for rename on name field and usage of regexp
		{
			Query: &Query{
				Start: "-4h",
				End:   "-2h",
				Filter: map[string]string{
					"name": "env.temp.*",
				}},
			Expect: `from(bucket:"mybucket") |> range(start:-4h,stop:-2h) |> filter(fn: (r) => r._measurement =~ /^env.temp.*$/)`,
		},
		{
			Query: &Query{
				Start: "-4h",
				End:   "-2h",
				Tail:  intptr(123),
				Filter: map[string]string{
					"name": "env.temp.*",
				}},
			Expect: `from(bucket:"mybucket") |> range(start:-4h,stop:-2h) |> filter(fn: (r) => r._measurement =~ /^env.temp.*$/) |> tail(n:123)`,
		},
	}

	for _, c := range testcases {
		s, err := buildFluxQuery("mybucket", c.Query)

		if c.ShouldFail {
			if err == nil {
				t.Fatal(fmt.Printf("expected error"))
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if s != c.Expect {
				t.Fatalf("flux query expected:\nexpect: %s\noutput: %s", s, c.Expect)
			}
		}
	}
}

func TestBuildFluxBadQuery(t *testing.T) {
	testcases := []*Query{
		{
			Start: "-4h",
			Filter: map[string]string{
				"name": "); drop bucket",
			},
		},
		{
			Start: "); danger",
		},
		{
			End: "); danger",
		},
	}

	for _, query := range testcases {
		_, err := buildFluxQuery("mybucket", query)
		if err == nil {
			t.Fatalf("expected error for %#v", query)
		}
	}
}
