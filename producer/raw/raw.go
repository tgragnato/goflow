package rawproducer

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"time"

	"github.com/tgragnato/goflow/decoders/netflow"
	"github.com/tgragnato/goflow/decoders/netflowlegacy"
	"github.com/tgragnato/goflow/decoders/sflow"
	"github.com/tgragnato/goflow/producer"
)

// Producer that keeps the same format
// as the original flow samples.
// This can be used for debugging (eg: getting NetFlow Option Templates)
type RawProducer struct {
}

// Raw message
type RawMessage struct {
	Message      interface{}    `json:"message"`
	Src          netip.AddrPort `json:"src"`
	TimeReceived time.Time      `json:"time_received"`
}

func (m RawMessage) MarshalJSON() ([]byte, error) {
	typeStr := "unknown"
	switch m.Message.(type) {
	case *netflowlegacy.PacketNetFlowV5:
		typeStr = "netflowv5"
	case *netflow.NFv9Packet:
		typeStr = "netflowv9"
	case *netflow.IPFIXPacket:
		typeStr = "ipfix"
	case *sflow.Packet:
		typeStr = "sflow"
	}

	tmpStruct := struct {
		Type         string          `json:"type"`
		Message      interface{}     `json:"message"`
		Src          *netip.AddrPort `json:"src"`
		TimeReceived *time.Time      `json:"time_received"`
	}{
		Type:         typeStr,
		Message:      m.Message,
		Src:          &m.Src,
		TimeReceived: &m.TimeReceived,
	}
	return json.Marshal(tmpStruct)
}

func (m RawMessage) MarshalText() ([]byte, error) {
	var msgContents []byte
	var err error
	if msg, ok := m.Message.(interface {
		MarshalText() ([]byte, error)
	}); ok {
		msgContents, err = msg.MarshalText()
	}
	return []byte(fmt.Sprintf("%s %s: %s", m.TimeReceived.String(), m.Src.String(), string(msgContents))), err
}

func (p *RawProducer) Produce(msg interface{}, args *producer.ProduceArgs) ([]producer.ProducerMessage, error) {
	// should return msg wrapped
	// []*interface{msg,}
	return []producer.ProducerMessage{RawMessage{msg, args.Src, args.TimeReceived}}, nil
}

func (p *RawProducer) Commit(flowMessageSet []producer.ProducerMessage) {}

func (p *RawProducer) Close() {}
