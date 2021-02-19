package main

import (
	"fmt"
	"strings"
)

// Query holds a query time range and filters.
type Query struct {
	Start  string            `json:"start,omitempty"`
	End    string            `json:"end,omitempty"`
	Filter map[string]string `json:"filter"`
}

func buildFluxQuery(bucket string, query *Query) (string, error) {
	var fluxParts []string

	fluxParts = append(fluxParts, fmt.Sprintf(`from(bucket:"%s")`, bucket))

	var rangeParts []string

	if query.Start != "" {
		rangeParts = append(rangeParts, "start:"+query.Start)
	}
	if query.End != "" {
		rangeParts = append(rangeParts, "end:"+query.End)
	}

	if len(rangeParts) > 0 {
		part := fmt.Sprintf(`range(%s)`, strings.Join(rangeParts, ","))
		fluxParts = append(fluxParts, part)
	}

	var filterParts []string

	// TODO sanitize filter parts
	for field, pattern := range query.Filter {
		// handle spcial rename cases
		if field == "name" {
			field = "_measurement"
		}
		// handle wildcard or exact match. (this may not actually be an optimization)
		if strings.Contains(pattern, "*") {
			filterParts = append(filterParts, fmt.Sprintf("r.%s =~ /^%s$/", field, pattern))
		} else {
			filterParts = append(filterParts, fmt.Sprintf("r.%s == \"%s\"", field, pattern))
		}
	}

	if len(filterParts) > 0 {
		part := fmt.Sprintf(`filter(fn: (r) => %s)`, strings.Join(filterParts, " and "))
		fluxParts = append(fluxParts, part)
	}

	return strings.Join(fluxParts, " |> "), nil
}
