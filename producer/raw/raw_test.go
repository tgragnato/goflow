package rawproducer_test

import (
	"fmt"
	"net/netip"
	"testing"
	"time"

	"tgragnato.it/goflow/decoders/netflow"
	"tgragnato.it/goflow/decoders/netflowlegacy"
	"tgragnato.it/goflow/decoders/sflow"
	"tgragnato.it/goflow/producer"
	rawproducer "tgragnato.it/goflow/producer/raw"
)

func TestRawProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		msg     any
		wantErr bool
	}{
		{
			name:    "netflow v9 packet",
			msg:     &netflow.NFv9Packet{Version: 9},
			wantErr: false,
		},
		{
			name:    "netflow v5 packet",
			msg:     &netflowlegacy.PacketNetFlowV5{Version: 5},
			wantErr: false,
		},
		{
			name:    "netflow ipfix packet",
			msg:     &netflow.IPFIXPacket{Version: 10},
			wantErr: false,
		},
		{
			name:    "sflow packet",
			msg:     &sflow.Packet{Version: 5},
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
			if !tt.wantErr {
				// Test MarshalText
				rawMsg, ok := msgs[0].(rawproducer.RawMessage)
				if !ok {
					t.Errorf("expected RawMessage, got %T", msgs[0])
				}
				text, err := rawMsg.MarshalText()
				if err != nil {
					t.Errorf("MarshalText() error = %v", err)
				}
				if len(text) == 0 {
					t.Errorf("MarshalText() returned empty bytes")
				}
				// Test MarshalJSON
				jsonBytes, err := rawMsg.MarshalJSON()
				if err != nil {
					t.Errorf("MarshalJSON() error = %v", err)
				}
				if len(jsonBytes) == 0 {
					t.Errorf("MarshalJSON() returned empty bytes")
				}
			}
		})
	}
}

func TestRawProducerCommitAndClose(t *testing.T) {
	t.Parallel()

	p := &rawproducer.RawProducer{}
	// Test Commit with empty slice
	p.Commit(nil)
	// Test Commit with messages
	args := &producer.ProduceArgs{
		Src:          netip.MustParseAddrPort("127.0.0.1:1234"),
		TimeReceived: time.Now(),
	}
	msgs, err := p.Produce(&netflow.NFv9Packet{Version: 9}, args)
	if err != nil {
		t.Fatalf("Produce() error = %v", err)
	}
	p.Commit(msgs)
	// Test Close
	p.Close()
}

func TestRawProducerMarshalTextError(t *testing.T) {
	t.Parallel()

	// Test MarshalText with a message that implements MarshalText but returns an error
	rawMsg := rawproducer.RawMessage{
		Message:      &errorMarshaler{},
		Src:          netip.MustParseAddrPort("127.0.0.1:1234"),
		TimeReceived: time.Now(),
	}
	_, err := rawMsg.MarshalText()
	if err == nil {
		t.Errorf("expected error for MarshalText with error-returning message")
	}
}

type errorMarshaler struct{}

func (e *errorMarshaler) MarshalText() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}
