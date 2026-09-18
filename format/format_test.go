package format_test

import (
	"testing"

	"tgragnato.it/goflow/format"
)

type testFormatDriver struct {
	name string
}

func (d *testFormatDriver) Prepare() error { return nil }
func (d *testFormatDriver) Init() error    { return nil }
func (d *testFormatDriver) Format(data any) ([]byte, []byte, error) {
	return []byte(d.name), []byte("text"), nil
}

func TestFormatDriver(t *testing.T) {
	t.Parallel()
	driver := &testFormatDriver{name: "test"}
	
	key, text, err := driver.Format("data")
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if string(key) != "test" {
		t.Fatalf("expected key 'test', got %q", string(key))
	}
	if string(text) != "text" {
		t.Fatalf("expected text 'text', got %q", string(text))
	}
}

func TestFormatInterface(t *testing.T) {
	t.Parallel()
	driver := &testFormatDriver{name: "test"}
	
	var _ format.FormatInterface = driver
	
	key, text, err := driver.Format("data")
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if string(key) != "test" || string(text) != "text" {
		t.Fatalf("unexpected result: key=%q text=%q", string(key), string(text))
	}
}

func TestDriverFormatError(t *testing.T) {
	t.Parallel()
	err := &format.DriverFormatError{
		Driver: "test",
		Err:    format.ErrFormat,
	}
	
	if err.Error() == "" {
		t.Fatal("DriverFormatError.Error() returned empty string")
	}
	
	unwrapped := err.Unwrap()
	if len(unwrapped) != 2 {
		t.Fatalf("expected 2 unwrapped errors, got %d", len(unwrapped))
	}
}

func TestFormatErrorVariables(t *testing.T) {
	t.Parallel()
	if format.ErrFormat == nil {
		t.Fatal("ErrFormat should not be nil")
	}
	if format.ErrNoSerializer == nil {
		t.Fatal("ErrNoSerializer should not be nil")
	}
}

func TestFormatStruct(t *testing.T) {
	t.Parallel()
	driver := &testFormatDriver{name: "test"}
	
	f := &format.Format{
		FormatDriver: driver,
	}
	
	key, text, err := f.Format("data")
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if string(key) != "test" {
		t.Fatalf("expected key 'test', got %q", string(key))
	}
	if string(text) != "text" {
		t.Fatalf("expected text 'text', got %q", string(text))
	}
}

func TestFormatStructError(t *testing.T) {
	t.Parallel()
	
	// Test with a driver that returns an error
	errDriver := &errorFormatDriver{}
	f := &format.Format{
		FormatDriver: errDriver,
	}
	
	_, _, err := f.Format("data")
	if err == nil {
		t.Fatal("expected error from Format")
	}
	
	// Check that it's a DriverFormatError
	if _, ok := err.(*format.DriverFormatError); !ok {
		t.Fatalf("expected DriverFormatError, got %T", err)
	}
}

type errorFormatDriver struct{}

func (d *errorFormatDriver) Prepare() error { return nil }
func (d *errorFormatDriver) Init() error    { return nil }
func (d *errorFormatDriver) Format(data any) ([]byte, []byte, error) {
	return nil, nil, format.ErrFormat
}
