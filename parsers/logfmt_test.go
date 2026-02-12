package parsers

import (
	"encoding/json"
	"testing"
)

func TestParseLogfmt(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantOK  bool
		wantObj map[string]any
	}{
		{
			name:   "basic key=value pairs",
			line:   `level=info msg="hello world" latency=1.23`,
			wantOK: true,
			wantObj: map[string]any{
				"level":   "info",
				"msg":     "hello world",
				"latency": 1.23,
			},
		},
		{
			name:   "integer and bool coercion",
			line:   `status=200 ok=true size=512`,
			wantOK: true,
			wantObj: map[string]any{
				"status": 200,
				"ok":     true,
				"size":   512,
			},
		},
		{
			name:   "false bool coercion",
			line:   `error=false`,
			wantOK: true,
			wantObj: map[string]any{
				"error": false,
			},
		},
		{
			name:   "empty value",
			line:   `key=`,
			wantOK: true,
			wantObj: map[string]any{
				"key": "",
			},
		},
		{
			name:   "bare words are valid logfmt keys with empty values",
			line:   "info started",
			wantOK: true,
			wantObj: map[string]any{
				"info":    "",
				"started": "",
			},
		},
		{
			name:   "empty line",
			line:   "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseLogfmt(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}

			var obj map[string]any
			if err := json.Unmarshal([]byte(got), &obj); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			for k, wantVal := range tt.wantObj {
				gotVal, exists := obj[k]
				if !exists {
					t.Errorf("missing key %q", k)
					continue
				}
				// json.Unmarshal always uses float64 for numbers
				switch w := wantVal.(type) {
				case int:
					if gotVal != float64(w) {
						t.Errorf("key %q: got %v, want %v", k, gotVal, w)
					}
				default:
					if gotVal != wantVal {
						t.Errorf("key %q: got %v (%T), want %v (%T)", k, gotVal, gotVal, wantVal, wantVal)
					}
				}
			}
		})
	}
}
