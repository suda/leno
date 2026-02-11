package parsers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/satyrius/gonx"
)

// nginxIngressFormat is the default log format used by nginx-ingress controller.
// It is a superset of the standard combined log format with upstream and request metadata.
// See: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/
const nginxIngressFormat = `$remote_addr - $remote_user [$time_local] "$request" ` +
	`$status $body_bytes_sent "$http_referer" "$http_user_agent" ` +
	`$request_length $request_time [$proxy_upstream_name] ` +
	`[$proxy_alternative_upstream_name] $upstream_addr ` +
	`$upstream_response_length $upstream_response_time ` +
	`$upstream_status $req_id`

// combinedLogFormat is the standard nginx combined log format.
const combinedLogFormat = `$remote_addr - $remote_user [$time_local] "$request" ` +
	`$status $body_bytes_sent "$http_referer" "$http_user_agent"`

var (
	nginxIngressParser = gonx.NewParser(nginxIngressFormat)
	combinedParser     = gonx.NewParser(combinedLogFormat)
)

// ParseNginx tries to parse a line as nginx-ingress extended format first,
// falling back to standard combined log format, and returns a JSON string.
// Returns ("", false) if the line doesn't match either format.
func ParseNginx(line string) (string, bool) {
	entry, err := nginxIngressParser.ParseString(line)
	if err != nil {
		entry, err = combinedParser.ParseString(line)
		if err != nil {
			return "", false
		}
	}

	obj := make(map[string]any)

	str := func(field string) string {
		v, _ := entry.Field(field)
		return v
	}

	setStr := func(key, field string) {
		if v := str(field); v != "" && v != "-" {
			obj[key] = v
		}
	}

	setInt := func(key, field string) {
		if v := str(field); v != "" && v != "-" {
			if n, err := strconv.Atoi(v); err == nil {
				obj[key] = n
			}
		}
	}

	setFloat := func(key, field string) {
		if v := str(field); v != "" && v != "-" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				obj[key] = f
			}
		}
	}

	setStr("remote_addr", "remote_addr")
	setStr("remote_user", "remote_user")
	setStr("time_local", "time_local")
	setStr("request", "request")
	setInt("status", "status")
	setInt("body_bytes", "body_bytes_sent")
	setStr("http_referer", "http_referer")
	setStr("http_user_agent", "http_user_agent")

	// Parse time_local into RFC3339 UTC
	if raw := str("time_local"); raw != "" {
		if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", raw); err == nil {
			obj["time"] = t.UTC().Format(time.RFC3339)
		}
	}

	// Split request into method / path / http_version
	if req := str("request"); req != "" {
		if parts := splitRequest(req); parts != nil {
			obj["method"] = parts[0]
			obj["path"] = parts[1]
			obj["http_version"] = parts[2]
		}
	}

	// nginx-ingress extended fields
	setInt("request_length", "request_length")
	setFloat("request_time", "request_time")
	setStr("upstream_name", "proxy_upstream_name")
	setStr("upstream_alt_name", "proxy_alternative_upstream_name")
	setStr("upstream_addr", "upstream_addr")
	setInt("upstream_response_length", "upstream_response_length")
	setStr("upstream_response_time", "upstream_response_time")
	setInt("upstream_status", "upstream_status")
	setStr("request_id", "req_id")

	b, err := json.Marshal(obj)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// splitRequest splits "METHOD /path HTTP/x.x" into its three components.
func splitRequest(req string) []string {
	parts := make([]string, 0, 3)
	i := 0
	for k := 0; k < 2; k++ {
		j := i
		for j < len(req) && req[j] != ' ' {
			j++
		}
		parts = append(parts, req[i:j])
		i = j + 1
	}
	if i <= len(req) {
		parts = append(parts, req[i:])
	}
	if len(parts) != 3 {
		return nil
	}
	return parts
}
