package syslog

import (
	"flag"
	"log"
	"log/syslog"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var syslogStateMu sync.Mutex

func TestSyslogDriverPrepare(t *testing.T) {
	t.Parallel()
	syslogStateMu.Lock()
	defer syslogStateMu.Unlock()

	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	d := &SyslogDriver{}
	if err := d.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if got := flag.Lookup("transport.syslog.protocol").Value.String(); got != "udp" {
		t.Fatalf("protocol = %q, want udp", got)
	}
	if got := flag.Lookup("transport.syslog.address").Value.String(); got != "localhost:514" {
		t.Fatalf("address = %q, want localhost:514", got)
	}
}

func TestSyslogDriverInitError(t *testing.T) {
	t.Parallel()
	d := &SyslogDriver{protocol: "udp", address: "invalid:0"}
	if err := d.Init(); err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestSyslogDriverInitSendAndClose(t *testing.T) {
	t.Parallel()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer conn.Close()

	remote, err := syslog.Dial("udp", conn.LocalAddr().String(), syslog.LOG_INFO|syslog.LOG_LOCAL0, "goflow")
	if err != nil {
		t.Fatalf("syslog.Dial failed: %v", err)
	}
	defer remote.Close()

	d := &SyslogDriver{protocol: "udp", address: conn.LocalAddr().String()}
	if err := d.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer log.SetOutput(log.Writer())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain message", "plain data", "plain data"},
		{"timestamp stripped", "2024/01/01 12:00:00 timestamped data", "timestamped data"},
		{"whitespace trimmed", "2024/01/01 12:00:00  spaced  data  ", "spaced  data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conn, err := net.ListenPacket("udp", "127.0.0.1:0")
			if err != nil {
				t.Fatalf("listen failed: %v", err)
			}
			defer conn.Close()

			remote, err := syslog.Dial("udp", conn.LocalAddr().String(), syslog.LOG_INFO|syslog.LOG_LOCAL0, "goflow")
			if err != nil {
				t.Fatalf("syslog.Dial failed: %v", err)
			}
			originalOutput := log.Writer()
			defer func() {
				syslogStateMu.Lock()
				log.SetOutput(originalOutput)
				syslogStateMu.Unlock()
				remote.Close()
			}()

			d := &SyslogDriver{protocol: "udp", address: conn.LocalAddr().String()}
			syslogStateMu.Lock()
			if err := d.Init(); err != nil {
				syslogStateMu.Unlock()
				t.Fatalf("Init failed: %v", err)
			}
			if err := d.Send([]byte("key"), []byte(tt.input)); err != nil {
				syslogStateMu.Unlock()
				t.Fatalf("Send failed: %v", err)
			}
			syslogStateMu.Unlock()

			conn.SetReadDeadline(time.Now().Add(time.Second))
			buf := make([]byte, 1024)
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				t.Fatalf("read syslog message: %v", err)
			}
			got := string(buf[:n])
			if !strings.Contains(got, tt.want) {
				t.Fatalf("syslog message = %q, want it to contain %q", got, tt.want)
			}
		})
	}

	if err := d.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
