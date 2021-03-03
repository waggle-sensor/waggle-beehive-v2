package main

import (
	"testing"
)

func TestBuildFluxQuery(t *testing.T) {
	testcases := []struct {
		Query  *Query
		Expect string
	}{
		// this case checks for start in time range
		{
			Query: &Query{
				Start: "-4h",
			},
			Expect: `from(bucket:"mybucket") |> range(start:-4h)`,
		},
		// this case checks for start and end in time range
		{
			Query: &Query{
				Start: "-4h",
				End:   "-2h",
			},
			Expect: `from(bucket:"mybucket") |> range(start:-4h,stop:-2h)`,
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
	}

	for _, c := range testcases {
		s, err := buildFluxQuery("mybucket", c.Query)

		if err != nil {
			t.Fatal(err)
		}

		if s != c.Expect {
			t.Fatalf("flux query expected:\nexpect: %s\noutput: %s", s, c.Expect)
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
