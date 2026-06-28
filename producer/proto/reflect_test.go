package protoproducer

import (
	"bytes"
	"testing"
)

func TestGetBytes(t *testing.T) {
	t.Parallel()
	d := []byte{0xAA, 0x55, 0xAB, 0x56}

	check := func(got, want []byte) {
		t.Helper()
		if !bytes.Equal(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
	checkNil := func(got []byte) {
		t.Helper()
		if got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	}

	check(GetBytes(d, 16, 16, true), []byte{0xAB, 0x56})
	check(GetBytes(d, 24, 8, true), []byte{0x56})
	check(GetBytes(d, 24, 32, true), []byte{0x56, 0x00, 0x00, 0x00})

	checkNil(GetBytes(d, 32, 0, true))

	check(GetBytes(d, 32, 16, true), []byte{0x00, 0x00})
	check(GetBytes(d, 4, 16, true), []byte{0xA5, 0x5A})
	check(GetBytes(d, 4, 16, false), []byte{0xA5, 0x5A})
	check(GetBytes(d, 4, 4, true), []byte{0x0A})
	check(GetBytes(d, 4, 4, false), []byte{0xA0})
	check(GetBytes(d, 4, 6, true), []byte{0x29})
	check(GetBytes(d, 4, 6, false), []byte{0xA4})
	check(GetBytes(d, 20, 6, true), []byte{0x2D})
	check(GetBytes(d, 20, 6, false), []byte{0xB4})
	check(GetBytes(d, 5, 10, true), []byte{0x4A, 0x02})
	check(GetBytes(d, 30, 10, true), []byte{0x80, 0x00})
	check(GetBytes(d, 30, 10, false), []byte{0x80, 0x00})
	check(GetBytes(d, 30, 2, true), []byte{0x02})
	check(GetBytes(d, 30, 2, false), []byte{0x80})
	check(GetBytes(d, 32, 1, true), []byte{0})
}

func BenchmarkGetBytes(b *testing.B) {
	d := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		GetBytes(d, 2, 10, false)
	}
}
