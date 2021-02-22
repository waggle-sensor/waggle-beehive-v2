package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Query holds a query time range and filters.
type Query struct {
	Start  string            `json:"start,omitempty"`
	End    string            `json:"end,omitempty"`
	Filter map[string]string `json:"filter"`
}

// buildFluxQuery builds a Flux query string for InfluxDB from a bucket name and Query
func buildFluxQuery(bucket string, query *Query) (string, error) {
	// start query out with data bucket
	parts := []string{
		fmt.Sprintf(`from(bucket:"%s")`, bucket),
	}

	// add range subquery, if not empty
	rangeSubquery, err := buildRangeSubquery(query)
	if err != nil {
		return "", err
	}
	if rangeSubquery != "" {
		parts = append(parts, rangeSubquery)
	}

	// add filter subquery, if not empty
	filterSubquery, err := buildFilterSubquery(query)
	if err != nil {
		return "", err
	}
	if filterSubquery != "" {
		parts = append(parts, filterSubquery)
	}

	return strings.Join(parts, " |> "), nil
}

func buildRangeSubquery(query *Query) (string, error) {
	var parts []string
	if !isValidFilterString(query.Start) {
		return "", fmt.Errorf("invalid start timestamp %q", query.Start)
	}
	if !isValidFilterString(query.End) {
		return "", fmt.Errorf("invalid end timestamp %q", query.End)
	}
	if query.Start != "" {
		parts = append(parts, "start:"+query.Start)
	}
	if query.End != "" {
		parts = append(parts, "end:"+query.End)
	}
	if len(parts) > 0 {
		return fmt.Sprintf(`range(%s)`, strings.Join(parts, ",")), nil
	}
	return "", nil
}

// Waggle and InfluxDB use slightly different field names in at least one case, so
// we keep a map to document and translate between them when needed.
var fieldRenameMap = map[string]string{
	"name": "_measurement",
}

func buildFilterSubquery(query *Query) (string, error) {
	var parts []string

	for field, pattern := range query.Filter {
		if !isValidFilterString(field) {
			return "", fmt.Errorf("invalid filter field name %q", field)
		}
		if !isValidFilterString(pattern) {
			return "", fmt.Errorf("invalid filter field pattern %q", pattern)
		}

		// rename field, if needed
		if s, ok := fieldRenameMap[field]; ok {
			field = s
		}

		// handle wildcard or exact match. (this may not actually be an optimization)
		if strings.Contains(pattern, "*") {
			parts = append(parts, fmt.Sprintf("r.%s =~ /^%s$/", field, pattern))
		} else {
			parts = append(parts, fmt.Sprintf("r.%s == \"%s\"", field, pattern))
		}
	}

	if len(parts) > 0 {
		return fmt.Sprintf(`filter(fn: (r) => %s)`, strings.Join(parts, " and ")), nil
	}

	return "", nil
}

var validQueryStringRE = regexp.MustCompile("^[A-Za-z0-9+-_.*: ]*$")

func isValidFilterString(s string) bool {
	return len(s) < 128 && validQueryStringRE.MatchString(s)
}
