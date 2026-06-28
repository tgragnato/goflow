package utils

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestUDPReceiver(t *testing.T) {
	t.Parallel()
	addr := "::1"
	port, err := getFreeUDPPort()
	if err != nil {
		t.Fatalf("getFreeUDPPort: %v", err)
	}
	t.Logf("starting UDP receiver on %s:%d\n", addr, port)

	r, err := NewUDPReceiver(nil)
	if err != nil {
		t.Fatalf("NewUDPReceiver: %v", err)
	}

	if err := r.Start(addr, port, nil); err != nil {
		t.Fatalf("Start: %v", err)
	}
	sendMessage := func(msg string) error {
		conn, err := net.Dial("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
		if err != nil {
			return fmt.Errorf("dial udp: %w", err)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			if closeErr := conn.Close(); closeErr != nil {
				return fmt.Errorf("close udp after write failure: %w", closeErr)
			}
			return fmt.Errorf("write udp: %w", err)
		}
		if err := conn.Close(); err != nil {
			return fmt.Errorf("close udp: %w", err)
		}
		return nil
	}
	if err := sendMessage("message"); err != nil {
		t.Fatalf("sendMessage: %v", err)
	}
	t.Log("sending message\n")
	if err := r.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestUDPClose(t *testing.T) {
	t.Parallel()
	addr := "::1"
	port, err := getFreeUDPPort()
	if err != nil {
		t.Fatalf("getFreeUDPPort: %v", err)
	}
	t.Logf("starting UDP receiver on %s:%d\n", addr, port)

	r, err := NewUDPReceiver(nil)
	if err != nil {
		t.Fatalf("NewUDPReceiver: %v", err)
	}
	if err := r.Start(addr, port, nil); err != nil {
		t.Fatalf("first Start: %v", err)
	}
	if err := r.Stop(); err != nil {
		t.Fatalf("first Stop: %v", err)
	}
	if err := r.Start(addr, port, nil); err != nil {
		t.Fatalf("second Start: %v", err)
	}
	if err := r.Start(addr, port, nil); err == nil {
		t.Fatal("expected error on third Start, got nil")
	}
	if err := r.Stop(); err != nil {
		t.Fatalf("second Stop: %v", err)
	}
	if err := r.Stop(); err == nil {
		t.Fatal("expected error on second Stop, got nil")
	}
}

func TestUDPReceiverDrainOnStop(t *testing.T) {
	t.Parallel()
	cfg := &UDPReceiverConfig{
		Workers:   1,
		Sockets:   1,
		QueueSize: 1000,
	}
	r, err := NewUDPReceiver(cfg)
	if err != nil {
		t.Fatalf("NewUDPReceiver: %v", err)
	}

	var decoded atomic.Int64
	decodeFunc := func(msg interface{}) error {
		decoded.Add(1)
		time.Sleep(2 * time.Millisecond) // slow decode to ensure backlog exists
		return nil
	}

	total := 50
	r.ready = make(chan bool) // mark as started without opening sockets
	if err := r.decoders(cfg.Workers, decodeFunc); err != nil {
		t.Fatalf("decoders: %v", err)
	}
	for i := 0; i < total; i++ {
		r.dispatch <- &udpPacket{
			src:      &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234},
			dst:      &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 5678},
			size:     1,
			payload:  []byte{1},
			received: time.Now().UTC(),
		}
	}

	if err := r.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if got := decoded.Load(); got != int64(total) {
		t.Fatalf("expected %d decoded, got %d", total, got)
	}
}

func TestUDPReceiverDecodeError(t *testing.T) {
	t.Parallel()
	addr := "::1"
	port, err := getFreeUDPPort()
	if err != nil {
		t.Fatalf("getFreeUDPPort: %v", err)
	}

	r, err := NewUDPReceiver(nil)
	if err != nil {
		t.Fatalf("NewUDPReceiver: %v", err)
	}

	wantErr := errors.New("decode error")
	errReady := make(chan struct{})
	gotErr := make(chan error, 1)
	go func() {
		close(errReady)
		err := <-r.Errors()
		gotErr <- err
	}()
	<-errReady

	decodeFunc := func(msg interface{}) error {
		return wantErr
	}

	if err := r.Start(addr, port, decodeFunc); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer func() {
		if err := r.Stop(); err != nil {
			t.Errorf("Stop: %v", err)
		}
	}()

	conn, err := net.Dial("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	if _, err = conn.Write([]byte("message")); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close conn: %v", err)
	}

	select {
	case err := <-gotErr:
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected %v, got %v", wantErr, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for decoder error")
	}
}

func getFreeUDPPort() (int, error) {
	a, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("resolve udp addr: %w", err)
	}
	l, err := net.ListenUDP("udp", a)
	if err != nil {
		return 0, fmt.Errorf("listen udp: %w", err)
	}
	port := l.LocalAddr().(*net.UDPAddr).Port
	if err := l.Close(); err != nil {
		return 0, fmt.Errorf("close udp listener: %w", err)
	}
	return port, nil
}
