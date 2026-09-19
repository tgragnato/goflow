package transport

import (
	"errors"
	"slices"
	"sync"
	"testing"
)

type stubDriver struct {
	prepareErr error
	initErr    error
	closeErr   error
	sendErr    error
	prepare    bool
	init       bool
	close      bool
	sendKey    []byte
	sendData   []byte
}

func (s *stubDriver) Prepare() error {
	s.prepare = true
	return s.prepareErr
}

func (s *stubDriver) Init() error {
	s.init = true
	return s.initErr
}

func (s *stubDriver) Close() error {
	s.close = true
	return s.closeErr
}

func (s *stubDriver) Send(key, data []byte) error {
	s.sendKey = key
	s.sendData = data
	return s.sendErr
}

var transportTestMu sync.Mutex

func TestRegisterTransportDriver(t *testing.T) {
	t.Parallel()
	transportTestMu.Lock()
	defer transportTestMu.Unlock()

	d := &stubDriver{}
	RegisterTransportDriver("stub", d)
	if !d.prepare {
		t.Fatal("Prepare was not called")
	}

	t.Run("panics on prepare error", func(t *testing.T) {
		t.Parallel()
		d2 := &stubDriver{prepareErr: errors.New("prepare failed")}
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic")
			}
		}()
		RegisterTransportDriver("stub2", d2)
	})
}

func TestFindTransport(t *testing.T) {
	t.Parallel()
	transportTestMu.Lock()
	defer transportTestMu.Unlock()

	d := &stubDriver{initErr: errors.New("init failed")}
	RegisterTransportDriver("find-test-stub", d)

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		_, err := FindTransport("missing")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("init error wrapped", func(t *testing.T) {
		t.Parallel()
		_, err := FindTransport("find-test-stub")
		if err == nil {
			t.Fatal("expected error")
		}
		var dte *DriverTransportError
		if !errors.As(err, &dte) {
			t.Fatalf("expected DriverTransportError, got %T", err)
		}
		if dte.Driver != "find-test-stub" {
			t.Fatalf("driver = %q, want find-test-stub", dte.Driver)
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		d2 := &stubDriver{}
		RegisterTransportDriver("find-test-success", d2)
		tr, err := FindTransport("find-test-success")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tr == nil {
			t.Fatal("expected transport")
		}
		if !d2.init {
			t.Fatal("Init was not called")
		}
	})
}

func TestTransportSend(t *testing.T) {
	t.Parallel()
	transportTestMu.Lock()
	defer transportTestMu.Unlock()

	d := &stubDriver{}
	RegisterTransportDriver("send-test-stub", d)
	tr, _ := FindTransport("send-test-stub")

	if err := tr.Send([]byte("k"), []byte("v")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(d.sendKey) != "k" || string(d.sendData) != "v" {
		t.Fatalf("send mismatch key=%q data=%q", d.sendKey, d.sendData)
	}

	d.sendErr = errors.New("send failed")
	err := tr.Send([]byte("k"), []byte("v"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTransportClose(t *testing.T) {
	t.Parallel()
	transportTestMu.Lock()
	defer transportTestMu.Unlock()

	d := &stubDriver{}
	RegisterTransportDriver("close-test-stub", d)
	tr, _ := FindTransport("close-test-stub")

	if err := tr.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.close {
		t.Fatal("Close was not called")
	}

	d.closeErr = errors.New("close failed")
	err := tr.Close()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetTransports(t *testing.T) {
	t.Parallel()
	transportTestMu.Lock()
	defer transportTestMu.Unlock()

	RegisterTransportDriver("get-transports-a", &stubDriver{})
	RegisterTransportDriver("get-transports-b", &stubDriver{})
	names := GetTransports()
	for _, name := range []string{"get-transports-a", "get-transports-b"} {
		if !slices.Contains(names, name) {
			t.Fatalf("GetTransports() = %v, want it to contain %q", names, name)
		}
	}
}
