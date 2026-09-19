package logging

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		level       string
		format      string
		wantErr     bool
		wantErrText string
	}{
		{name: "info text", level: "info", format: "", wantErr: false},
		{name: "debug json", level: "debug", format: "json", wantErr: false},
		{name: "warn text", level: "warn", format: "text", wantErr: false},
		{name: "invalid level", level: "nope", format: "", wantErr: true, wantErrText: "parse log level"},
		{name: "empty level", level: "", format: "", wantErr: true, wantErrText: "parse log level"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger, err := NewLogger(tt.level, tt.format)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantErrText != "" && !contains(err.Error(), tt.wantErrText) {
					t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if logger == nil {
				t.Fatal("expected non-nil logger")
			}
		})
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
