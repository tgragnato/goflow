package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsNamespace(t *testing.T) {
	t.Parallel()

	if NAMESPACE != "goflow" {
		t.Errorf("expected namespace 'goflow', got %s", NAMESPACE)
	}
}

func TestMetricDescriptors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vec  prometheus.Collector
	}{
		{"MetricReceivedDroppedPackets", MetricReceivedDroppedPackets},
		{"MetricReceivedDroppedBytes", MetricReceivedDroppedBytes},
		{"MetricTrafficBytes", MetricTrafficBytes},
		{"MetricTrafficPackets", MetricTrafficPackets},
		{"DecoderErrors", DecoderErrors},
		{"NetFlowStats", NetFlowStats},
		{"NetFlowErrors", NetFlowErrors},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reg := prometheus.NewRegistry()
			if err := reg.Register(tt.vec); err != nil {
				t.Errorf("%s: register error: %v", tt.name, err)
			}
		})
	}
}

func TestMetricTrafficBytesIncrement(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	if err := reg.Register(MetricTrafficBytes); err != nil {
		t.Fatalf("register: %v", err)
	}

	MetricTrafficBytes.With(prometheus.Labels{
		"remote_ip":  "10.0.0.1",
		"local_ip":   "10.0.0.2",
		"local_port": "9999",
		"type":       "test",
	}).Add(42)

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}

	found := false
	for _, f := range families {
		if f.GetName() == "goflow_flow_traffic_bytes_total" {
			found = true
			for _, m := range f.GetMetric() {
				if m.GetCounter().GetValue() == 42 {
					return // success
				}
			}
		}
	}
	if !found {
		t.Errorf("expected metric family goflow_flow_traffic_bytes_total not found, got: %v", families)
	}
}

func TestPromDecoderWrapperReturnsErrorForWrongType(t *testing.T) {
	t.Parallel()

	wrapped := func(msg interface{}) error {
		return nil
	}

	fn := PromDecoderWrapper(wrapped, "test")
	err := fn("not a message")
	if err == nil {
		t.Errorf("expected error for wrong type, got nil")
	}
}
