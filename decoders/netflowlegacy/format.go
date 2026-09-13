package netflowlegacy

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// MarshalJSON encodes the packet without triggering MarshalText.
func (p *PacketNetFlowV5) MarshalJSON() ([]byte, error) {
	return json.Marshal(*p) // this is a trick to avoid having the JSON marshaller defaults to MarshalText
}

// MarshalText formats a concise text summary of the packet.
func (p *PacketNetFlowV5) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "NetFlowV%d seq:%d count:%d", p.Version, p.FlowSequence, p.Count), nil
}

// String renders a multi-line representation of the packet.
func (p PacketNetFlowV5) String() string {
	var str strings.Builder
	str.WriteString("NetFlow v5 Packet\n")
	str.WriteString("-----------------\n")
	fmt.Fprintf(&str, "  Version: %v\n", p.Version)
	fmt.Fprintf(&str, "  Count:  %v\n", p.Count)

	unixSeconds := time.Unix(int64(p.UnixSecs), int64(p.UnixNSecs))
	fmt.Fprintf(&str, "  SystemUptime: %v\n", time.Duration(p.SysUptime)*time.Millisecond)
	fmt.Fprintf(&str, "  UnixSeconds: %v\n", unixSeconds.String())
	fmt.Fprintf(&str, "  FlowSequence: %v\n", p.FlowSequence)
	fmt.Fprintf(&str, "  EngineType: %v\n", p.EngineType)
	fmt.Fprintf(&str, "  EngineId: %v\n", p.EngineId)
	fmt.Fprintf(&str, "  SamplingInterval: %v\n", p.SamplingInterval)
	fmt.Fprintf(&str, "  Records (%v):\n", len(p.Records))

	for i, record := range p.Records {
		fmt.Fprintf(&str, "    Record %v:\n", i)
		str.WriteString(record.String())
	}
	return str.String()
}

// String renders a multi-line representation of the record.
func (r RecordsNetFlowV5) String() string {
	str := fmt.Sprintf("      SrcAddr: %v\n", r.SrcAddr)
	str += fmt.Sprintf("      DstAddr: %v\n", r.DstAddr)
	str += fmt.Sprintf("      NextHop: %v\n", r.NextHop)
	str += fmt.Sprintf("      Input: %v\n", r.Input)
	str += fmt.Sprintf("      Output: %v\n", r.Output)
	str += fmt.Sprintf("      DPkts: %v\n", r.DPkts)
	str += fmt.Sprintf("      DOctets: %v\n", r.DOctets)
	str += fmt.Sprintf("      First: %v\n", time.Duration(r.First)*time.Millisecond)
	str += fmt.Sprintf("      Last: %v\n", time.Duration(r.Last)*time.Millisecond)
	str += fmt.Sprintf("      SrcPort: %v\n", r.SrcPort)
	str += fmt.Sprintf("      DstPort: %v\n", r.DstPort)
	str += fmt.Sprintf("      TCPFlags: %v\n", r.TCPFlags)
	str += fmt.Sprintf("      Proto: %v\n", r.Proto)
	str += fmt.Sprintf("      Tos: %v\n", r.Tos)
	str += fmt.Sprintf("      SrcAS: %v\n", r.SrcAS)
	str += fmt.Sprintf("      DstAS: %v\n", r.DstAS)
	str += fmt.Sprintf("      SrcMask: %v\n", r.SrcMask)
	str += fmt.Sprintf("      DstMask: %v\n", r.DstMask)

	return str
}
