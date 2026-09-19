package file

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestPrepare(t *testing.T) {
	t.Parallel()
	d := &FileDriver{lock: &sync.RWMutex{}}
	// Reset flag state before test to avoid redefinition
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := d.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
}

func TestInitStdout(t *testing.T) {
	t.Parallel()
	d := &FileDriver{fileDestination: "", lock: &sync.RWMutex{}}
	if err := d.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if d.w == nil {
		t.Fatal("stdout writer not initialized")
	}
	if d.w != os.Stdout {
		t.Fatal("expected os.Stdout writer")
	}
}

func TestInitFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")
	d := &FileDriver{fileDestination: path, lock: &sync.RWMutex{}}
	if err := d.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if d.file == nil {
		t.Fatal("file not opened")
	}
	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestInitFileError(t *testing.T) {
	t.Parallel()
	d := &FileDriver{fileDestination: "/nonexistent/dir/out.log", lock: &sync.RWMutex{}}
	if err := d.Init(); err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestSend(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	d := &FileDriver{w: buf, lineSeparator: "\n", lock: &sync.RWMutex{}}
	if err := d.Send([]byte("key"), []byte("data")); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if buf.String() != "data\n" {
		t.Fatalf("got %q, want %q", buf.String(), "data\n")
	}
}

func TestSendEmptySeparator(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	d := &FileDriver{w: buf, lineSeparator: "", lock: &sync.RWMutex{}}
	if err := d.Send([]byte("key"), []byte("data")); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if buf.String() != "data" {
		t.Fatalf("got %q, want %q", buf.String(), "data")
	}
}

func TestSendEmptyData(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	d := &FileDriver{w: buf, lineSeparator: "\n", lock: &sync.RWMutex{}}
	if err := d.Send([]byte("key"), nil); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if buf.String() != "\n" {
		t.Fatalf("got %q, want %q", buf.String(), "\n")
	}
}

func TestSendError(t *testing.T) {
	t.Parallel()
	errWriter := &errWriter{}
	d := &FileDriver{w: errWriter, lineSeparator: "\n", lock: &sync.RWMutex{}}
	if err := d.Send([]byte("key"), []byte("data")); err == nil {
		t.Fatal("expected error")
	}
}

type errWriter struct{}

func (w *errWriter) Write(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestClose(t *testing.T) {
	t.Parallel()
	d := &FileDriver{closed: false, lock: &sync.RWMutex{}}
	if err := d.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if !d.closed {
		t.Fatal("closed flag not set")
	}
}

func TestCloseAlreadyClosed(t *testing.T) {
	t.Parallel()
	d := &FileDriver{closed: true, lock: &sync.RWMutex{}}
	if err := d.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestCloseFileError(t *testing.T) {
	t.Parallel()
	file, err := os.CreateTemp("", "goflow-close-*")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	d := &FileDriver{
		fileDestination: file.Name(),
		file:            file,
		closed:          false,
		lock:            &sync.RWMutex{},
	}
	if err := d.Close(); err == nil {
		t.Fatal("expected error")
	}
}
