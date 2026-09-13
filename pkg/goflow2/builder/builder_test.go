package builder_test

import (
	"testing"

	_ "tgragnato.it/goflow/format/binary"
	_ "tgragnato.it/goflow/format/json"
	_ "tgragnato.it/goflow/format/text"
	"tgragnato.it/goflow/pkg/goflow2/builder"
	_ "tgragnato.it/goflow/transport/file"
	_ "tgragnato.it/goflow/transport/syslog"
)

func TestBuildFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "json formatter",
			input:   "json",
			wantErr: false,
		},
		{
			name:    "text formatter",
			input:   "text",
			wantErr: false,
		},
		{
			name:    "binary formatter",
			input:   "bin",
			wantErr: false,
		},
		{
			name:    "unknown formatter",
			input:   "unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := builder.BuildFormatter(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildFormatter(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBuildTransport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "file transport",
			input:   "file",
			wantErr: false,
		},
		{
			name:    "unknown transport",
			input:   "unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := builder.BuildTransport(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildTransport(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
