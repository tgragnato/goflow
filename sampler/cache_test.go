package sampler_test

import (
	"net"
	"testing"

	"tgragnato.it/goflow/sampler"
)

func TestInitAndGetHostnameByByteSlice(t *testing.T) {
	t.Parallel()

	sampler.Init()

	if got := sampler.GetHostnameByByteSlice(nil); got != "" {
		t.Fatalf("expected empty hostname for nil input, got %q", got)
	}

	// Exercise the public API with a valid IP without depending on DNS results.
	_ = sampler.GetHostnameByByteSlice(net.ParseIP("8.8.8.8"))
}
