package json_test

import (
	"bytes"
	"errors"
	"testing"

	jsonformat "tgragnato.it/goflow/format/json"
)

// mockJSONMarshaler implements json.Marshaler for testing.
type mockJSONMarshaler struct {
	data []byte
	err  error
}

func (m *mockJSONMarshaler) MarshalJSON() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

// mockKeyMarshalerJSON implements both Key() and json.Marshaler for testing.
type mockKeyMarshalerJSON struct {
	key  []byte
	data []byte
}

func (m *mockKeyMarshalerJSON) Key() []byte {
	return m.key
}

func (m *mockKeyMarshalerJSON) MarshalJSON() ([]byte, error) {
	return m.data, nil
}

func TestJsonDriverPrepare(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	if err := d.Prepare(); err != nil {
		t.Fatalf("Prepare returned error: %v", err)
	}
}

func TestJsonDriverInit(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	if err := d.Init(); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
}

func TestJsonFormatSuccess(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	data := map[string]any{"name": "test", "value": 42}

	key, payload, err := d.Format(data)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatal("expected non-empty payload")
	}
	if key != nil {
		t.Errorf("expected nil key, got %v", key)
	}
}

func TestJsonFormatWithKey(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	// A type that implements Key() but json.Marshal succeeds
	data := &mockKeyMarshalerJSON{
		key:  []byte("mykey"),
		data: []byte(`{"type":"custom"}`),
	}

	key, payload, err := d.Format(data)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if !bytes.Equal(key, []byte("mykey")) {
		t.Errorf("expected key 'mykey', got '%s'", key)
	}
	if len(payload) == 0 {
		t.Error("expected non-empty payload")
	}
}

func TestJsonFormatWithKeyAndMarshaler(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	data := &mockKeyMarshalerJSON{
		key:  []byte("mykey"),
		data: []byte(`{"type":"custom"}`),
	}

	key, payload, err := d.Format(data)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if !bytes.Equal(key, []byte("mykey")) {
		t.Errorf("expected key 'mykey', got '%s'", key)
	}
	if !bytes.Equal(payload, []byte(`{"type":"custom"}`)) {
		t.Errorf("unexpected payload: %s", payload)
	}
}

func TestJsonFormatMarshalError(t *testing.T) {
	t.Parallel()

	d := &jsonformat.JsonDriver{}
	data := &mockJSONMarshaler{err: errors.New("marshal error")}

	_, _, err := d.Format(data)
	if err == nil {
		t.Fatal("expected error for marshal failure")
	}
}
