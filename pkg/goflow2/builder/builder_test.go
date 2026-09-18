package builder_test

import (
	"testing"

	_ "tgragnato.it/goflow/format/binary"
	_ "tgragnato.it/goflow/format/json"
	_ "tgragnato.it/goflow/format/text"
	"tgragnato.it/goflow/pkg/goflow2/builder"
	"tgragnato.it/goflow/pkg/goflow2/config"
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
		{
			name:    "empty formatter",
			input:   "",
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
		{
			name:    "empty transport",
			input:   "",
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

func TestBuildProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "raw producer",
			cfg:     &config.Config{Produce: "raw"},
			wantErr: false,
		},
		{
			name:    "unknown produce value",
			cfg:     &config.Config{Produce: "unknown"},
			wantErr: true,
		},
		{
			name:    "empty produce value",
			cfg:     &config.Config{Produce: ""},
			wantErr: true,
		},
		{
			name:    "sample producer with missing mapping file",
			cfg:     &config.Config{Produce: "sample", MappingFile: "/nonexistent/mapping.yaml"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := builder.BuildProducer(tt.cfg, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildProducer(cfg) error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
