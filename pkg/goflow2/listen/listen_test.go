package listen_test

import (
	"testing"

	"tgragnato.it/goflow/pkg/goflow2/listen"
)

func TestParseListenAddressesBasic(t *testing.T) {
	t.Parallel()

	cfgs, err := listen.ParseListenAddresses("udp://0.0.0.0:2055")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(cfgs))
	}
	if cfgs[0].Scheme != "udp" {
		t.Errorf("expected scheme udp, got %s", cfgs[0].Scheme)
	}
	if cfgs[0].Port != 2055 {
		t.Errorf("expected port 2055, got %d", cfgs[0].Port)
	}
}

func TestParseListenAddressesMultiple(t *testing.T) {
	t.Parallel()

	cfgs, err := listen.ParseListenAddresses("udp://0.0.0.0:2055,tcp://127.0.0.1:9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(cfgs))
	}
}

func TestParseListenAddressesWithCount(t *testing.T) {
	t.Parallel()

	cfgs, err := listen.ParseListenAddresses("udp://0.0.0.0:2055?count=4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfgs[0].NumSockets != 4 {
		t.Errorf("expected 4 sockets, got %d", cfgs[0].NumSockets)
	}
}

func TestParseListenAddressesWithWorkers(t *testing.T) {
	t.Parallel()

	cfgs, err := listen.ParseListenAddresses("udp://0.0.0.0:2055?workers=8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfgs[0].NumWorkers != 8 {
		t.Errorf("expected 8 workers, got %d", cfgs[0].NumWorkers)
	}
}

func TestParseListenAddressesWithBlocking(t *testing.T) {
	t.Parallel()

	cfgs, err := listen.ParseListenAddresses("udp://0.0.0.0:2055?blocking=true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfgs[0].Blocking {
		t.Errorf("expected blocking true")
	}
}

func TestParseListenAddressesInvalidPort(t *testing.T) {
	t.Parallel()

	_, err := listen.ParseListenAddresses("udp://0.0.0.0:abc")
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestParseListenAddressesInvalidCount(t *testing.T) {
	t.Parallel()

	_, err := listen.ParseListenAddresses("udp://0.0.0.0:2055?count=bad")
	if err == nil {
		t.Fatal("expected error for invalid count")
	}
}
