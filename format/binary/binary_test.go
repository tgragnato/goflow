package binary_test

import (
	"bytes"
	"testing"

	"tgragnato.it/goflow/format/binary"
)

// mockKeyMarshaler implements both Key() and BinaryMarshaler for testing.
type mockKeyMarshaler struct {
	key  []byte
	data []byte
}

func (m *mockKeyMarshaler) Key() []byte {
	return m.key
}

func (m *mockKeyMarshaler) MarshalBinary() ([]byte, error) {
	return m.data, nil
}

// mockMarshaler implements encoding.BinaryMarshaler for testing.
type mockMarshaler struct {
	data []byte
	err  error
}

func (m *mockMarshaler) MarshalBinary() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

// mockKey implements Key() []byte for testing.
type mockKey struct {
	key []byte
}

func (m *mockKey) Key() []byte {
	return m.key
}

func TestBinaryDriverPrepare(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	if err := d.Prepare(); err != nil {
		t.Fatalf("Prepare returned error: %v", err)
	}
}

func TestBinaryDriverInit(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	if err := d.Init(); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
}

func TestBinaryFormatSuccess(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	data := &mockMarshaler{data: []byte("test data")}

	key, payload, err := d.Format(data)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if !bytes.Equal(payload, []byte("test data")) {
		t.Errorf("expected payload 'test data', got '%s'", payload)
	}
	if key != nil {
		t.Errorf("expected nil key, got %v", key)
	}
}

func TestBinaryFormatWithKey(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	data := &mockKey{key: []byte("mykey")}
	// This type doesn't implement BinaryMarshaler, so it should return ErrNoSerializer
	// But it does implement Key() so key should be extracted

	_, _, err := d.Format(data)
	if err == nil {
		t.Fatal("expected error for non-serializable type")
	}
}

func TestBinaryFormatWithKeyAndMarshaler(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	// A type that implements both Key() and BinaryMarshaler
	data := &mockKeyMarshaler{
		key:  []byte("mykey"),
		data: []byte("payload"),
	}

	key, payload, err := d.Format(data)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if !bytes.Equal(key, []byte("mykey")) {
		t.Errorf("expected key 'mykey', got '%s'", key)
	}
	if !bytes.Equal(payload, []byte("payload")) {
		t.Errorf("expected payload 'payload', got '%s'", payload)
	}
}

func TestBinaryFormatMarshalError(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	data := &mockMarshaler{err: bytes.ErrTooLarge}

	_, _, err := d.Format(data)
	if err == nil {
		t.Fatal("expected error for marshal failure")
	}
}

func TestBinaryFormatNoSerializer(t *testing.T) {
	t.Parallel()

	d := &binary.BinaryDriver{}
	// Simple type that implements neither Key nor BinaryMarshaler
	data := "plain string"

	_, _, err := d.Format(data)
	if err == nil {
		t.Fatal("expected error for non-serializable type")
	}
}
