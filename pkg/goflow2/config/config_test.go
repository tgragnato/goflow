package config_test

import (
	"flag"
	"testing"

	"tgragnato.it/goflow/pkg/goflow2/config"
)

func TestBindFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		args          []string
		wantListen    string
		wantProduce   string
		wantFormat    string
		wantTransport string
	}{
		{
			name:          "defaults",
			args:          []string{},
			wantListen:    "sflow://:6343,netflow://:2055",
			wantProduce:   "sample",
			wantFormat:    "json",
			wantTransport: "file",
		},
		{
			name:          "custom listen",
			args:          []string{"-listen", "netflow://:9999"},
			wantListen:    "netflow://:9999",
			wantProduce:   "sample",
			wantFormat:    "json",
			wantTransport: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			cfg := config.BindFlags(fs)
			if err := fs.Parse(tt.args); err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if cfg.ListenAddresses != tt.wantListen {
				t.Errorf("ListenAddresses = %q, want %q", cfg.ListenAddresses, tt.wantListen)
			}
			if cfg.Produce != tt.wantProduce {
				t.Errorf("Produce = %q, want %q", cfg.Produce, tt.wantProduce)
			}
			if cfg.Format != tt.wantFormat {
				t.Errorf("Format = %q, want %q", cfg.Format, tt.wantFormat)
			}
			if cfg.Transport != tt.wantTransport {
				t.Errorf("Transport = %q, want %q", cfg.Transport, tt.wantTransport)
			}
		})
	}
}
