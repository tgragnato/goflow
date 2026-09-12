package text_test

import (
	"bytes"
	"encoding"
	"testing"

	"tgragnato.it/goflow/format/text"
)

type mockTextMarshaler struct{}

func (m mockTextMarshaler) MarshalText() ([]byte, error) {
	return []byte("hello"), nil
}

type mockStringer struct{}

func (m mockStringer) String() string {
	return "world"
}

type mockKey struct{}

func (m mockKey) Key() []byte {
	return []byte("key")
}

func TestTextDriverFormat(t *testing.T) {
	t.Parallel()

	d := &text.TextDriver{}

	tests := []struct {
		name     string
		data     any
		wantKey  []byte
		wantData []byte
		wantErr  bool
	}{
		{
			name:     "text marshaler",
			data:     mockTextMarshaler{},
			wantKey:  nil,
			wantData: []byte("hello"),
			wantErr:  false,
		},
		{
			name:     "stringer",
			data:     mockStringer{},
			wantKey:  nil,
			wantData: []byte("world"),
			wantErr:  false,
		},
		{
			name: "key and text marshaler",
			data: struct {
				mockTextMarshaler
				mockKey
			}{},
			wantKey:  []byte("key"),
			wantData: []byte("hello"),
			wantErr:  false,
		},
		{
			name:     "unsupported",
			data:     42,
			wantKey:  nil,
			wantData: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			key, data, err := d.Format(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(key, tt.wantKey) {
				t.Errorf("key = %v, want %v", key, tt.wantKey)
			}
			if !bytes.Equal(data, tt.wantData) {
				t.Errorf("data = %v, want %v", data, tt.wantData)
			}
		})
	}
}

func TestTextDriverInitPrepare(t *testing.T) {
	t.Parallel()
	d := &text.TextDriver{}
	if err := d.Prepare(); err != nil {
		t.Errorf("Prepare() error = %v", err)
	}
	if err := d.Init(); err != nil {
		t.Errorf("Init() error = %v", err)
	}
}

// Ensure interface satisfaction at compile time.
var _ encoding.TextMarshaler = mockTextMarshaler{}
