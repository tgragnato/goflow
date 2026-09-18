package app

import (
	"context"
	"testing"
	"time"

	_ "tgragnato.it/goflow/format/binary"
	_ "tgragnato.it/goflow/format/json"
	_ "tgragnato.it/goflow/format/text"
	_ "tgragnato.it/goflow/transport/file"
	"tgragnato.it/goflow/pkg/goflow2/config"
)

func newTestConfig() *config.Config {
	return &config.Config{
		LogLevel:              "info",
		LogFmt:                "text",
		Format:                "bin",
		Transport:             "file",
		Produce:               "sample",
		ListenAddresses:       "sflow://127.0.0.1:6343",
		StoreJSONPath:         "/tmp/goflow-test-store",
		StoreJSONInterval:     time.Second,
		TemplatesTTL:          time.Minute,
		TemplatesExtendOnAccess: true,
		TemplatesSweepInterval:  time.Minute,
		SamplingRatesTTL:      time.Minute,
		SamplingRatesExtendOnAccess: true,
		SamplingRatesSweepInterval: time.Minute,
		ErrCnt:                10,
		ErrInt:                time.Millisecond,
	}
}

func TestAppNew(t *testing.T) {
	t.Parallel()
	app, err := New(newTestConfig())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if app == nil {
		t.Fatal("New() returned nil app")
	}
}

func TestAppNewWithEmptyAddr(t *testing.T) {
	t.Parallel()
	app, err := New(newTestConfig())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if app == nil {
		t.Fatal("New() returned nil app")
	}
}

func TestAppWait(t *testing.T) {
	t.Parallel()
	app, err := New(newTestConfig())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	
	ch := app.Wait()
	if ch == nil {
		t.Fatal("Wait() returned nil channel")
	}
}

func TestAppShutdown(t *testing.T) {
	t.Parallel()
	app, err := New(newTestConfig())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	app.Shutdown(ctx)
}

func TestAppRun(t *testing.T) {
	t.Parallel()
	app, err := New(newTestConfig())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	err = app.Run(ctx)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
}
