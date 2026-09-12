package rawproducer_test

import (
	"net/netip"
	"testing"
	"time"

	"tgragnato.it/goflow/decoders/netflow"
	"tgragnato.it/goflow/producer"
	rawproducer "tgragnato.it/goflow/producer/raw"
)

func TestRawProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		msg     interface{}
		wantErr bool
	}{
		{
			name:    "netflow v9 packet",
			msg:     &netflow.NFv9Packet{Version: 9},
			wantErr: false,
		},
		{
			name:    "netflow v5 packet",
			msg:     &netflow.NFv9Packet{Version: 5},
			wantErr: false,
		},
	}

	p := &rawproducer.RawProducer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args := &producer.ProduceArgs{
				Src:          netip.MustParseAddrPort("127.0.0.1:1234"),
				TimeReceived: time.Now(),
			}
			msgs, err := p.Produce(tt.msg, args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Produce() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(msgs) != 1 {
				t.Errorf("expected 1 message, got %d", len(msgs))
			}
		})
	}
}
