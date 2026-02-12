package parsers

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/go-logfmt/logfmt"
)

// ParseLogfmt parses a logfmt line and returns a JSON string.
// Returns ("", false) if the line contains no valid key=value pairs.
func ParseLogfmt(line string) (string, bool) {
	dec := logfmt.NewDecoder(bytes.NewBufferString(line))
	obj := make(map[string]any)

	for dec.ScanRecord() {
		for dec.ScanKeyval() {
			key := string(dec.Key())
			val := string(dec.Value())
			obj[key] = coerce(val)
		}
	}
	if err := dec.Err(); err != nil || len(obj) == 0 {
		return "", false
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// coerce tries to convert a string value to a number, bool, or keeps it as a string.
func coerce(v string) any {
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}
	return v
}
