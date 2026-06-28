package sflow

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tgragnato/goflow/decoders/utils"
)

func TestEncodeDecodeSFlow(t *testing.T) {
	t.Parallel()
	packet := Packet{
		Version:        5,
		IPVersion:      1,
		AgentIP:        utils.IPAddress{192, 0, 2, 1},
		SubAgentId:     1,
		SequenceNumber: 2,
		Uptime:         3,
		Samples: []interface{}{
			FlowSample{
				Header: SampleHeader{
					Format:               SAMPLE_FORMAT_FLOW,
					SampleSequenceNumber: 42,
					SourceIdType:         0,
					SourceIdValue:        7,
				},
				SamplingRate:     10,
				SamplePool:       20,
				Drops:            0,
				Input:            1,
				Output:           2,
				FlowRecordsCount: 1,
				Records: []FlowRecord{
					{
						Header: RecordHeader{
							DataFormat: FLOW_TYPE_RAW,
						},
						Data: SampledHeader{
							Protocol:       1,
							FrameLength:    64,
							Stripped:       0,
							OriginalLength: 4,
							HeaderData:     []byte{0xde, 0xad, 0xbe, 0xef},
						},
					},
				},
			},
		},
	}

	encoded, err := EncodeMessage(&packet)
	if err != nil {
		t.Fatalf("EncodeMessage: %v", err)
	}

	var decoded Packet
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), &decoded); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}
	if decoded.Version != uint32(5) {
		t.Fatalf("expected Version 5, got %d", decoded.Version)
	}
	if decoded.IPVersion != uint32(1) {
		t.Fatalf("expected IPVersion 1, got %d", decoded.IPVersion)
	}
	if !bytes.Equal(decoded.AgentIP, utils.IPAddress{192, 0, 2, 1}) {
		t.Fatalf("expected AgentIP {192,0,2,1}, got %v", decoded.AgentIP)
	}
	if len(decoded.Samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(decoded.Samples))
	}

	sample, ok := decoded.Samples[0].(FlowSample)
	if !ok {
		t.Fatal("expected FlowSample")
	}
	if sample.Header.SampleSequenceNumber != uint32(42) {
		t.Fatalf("expected SampleSequenceNumber 42, got %d", sample.Header.SampleSequenceNumber)
	}
	if sample.Header.SourceIdValue != uint32(7) {
		t.Fatalf("expected SourceIdValue 7, got %d", sample.Header.SourceIdValue)
	}
	if len(sample.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(sample.Records))
	}

	record := sample.Records[0]
	if record.Header.DataFormat != uint32(FLOW_TYPE_RAW) {
		t.Fatalf("expected DataFormat %d, got %d", FLOW_TYPE_RAW, record.Header.DataFormat)
	}
	header, ok := record.Data.(SampledHeader)
	if !ok {
		t.Fatal("expected SampledHeader")
	}
	if !bytes.Equal(header.HeaderData, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("expected HeaderData %v, got %v", []byte{0xde, 0xad, 0xbe, 0xef}, header.HeaderData)
	}
}

func TestEncodeDecodeSFlowExpandedFlowSample(t *testing.T) {
	t.Parallel()
	packet := Packet{
		Version:        5,
		IPVersion:      2,
		AgentIP:        utils.IPAddress{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		SubAgentId:     10,
		SequenceNumber: 11,
		Uptime:         12,
		Samples: []interface{}{
			ExpandedFlowSample{
				Header: SampleHeader{
					Format:               SAMPLE_FORMAT_EXPANDED_FLOW,
					SampleSequenceNumber: 100,
					SourceIdType:         2,
					SourceIdValue:        99,
				},
				SamplingRate:     400,
				SamplePool:       500,
				Drops:            2,
				InputIfFormat:    1,
				InputIfValue:     10,
				OutputIfFormat:   2,
				OutputIfValue:    20,
				FlowRecordsCount: 3,
				Records: []FlowRecord{
					{
						Data: SampledIPv4{
							SampledIPBase: SampledIPBase{
								Length:   60,
								Protocol: 6,
								SrcIP:    utils.IPAddress{192, 0, 2, 10},
								DstIP:    utils.IPAddress{198, 51, 100, 20},
								SrcPort:  12345,
								DstPort:  443,
								TcpFlags: 0x12,
							},
							Tos: 16,
						},
					},
					{
						Data: ExtendedSwitch{
							SrcVlan:     100,
							SrcPriority: 1,
							DstVlan:     200,
							DstPriority: 2,
						},
					},
					{
						Data: ExtendedGateway{
							NextHopIPVersion:  1,
							NextHop:           utils.IPAddress{203, 0, 113, 1},
							AS:                64512,
							SrcAS:             64513,
							SrcPeerAS:         64514,
							ASDestinations:    1,
							ASPathType:        1,
							ASPathLength:      2,
							ASPath:            []uint32{64515, 64516},
							CommunitiesLength: 2,
							Communities:       []uint32{100, 200},
							LocalPref:         300,
						},
					},
				},
			},
		},
	}

	encoded, err := EncodeMessage(&packet)
	if err != nil {
		t.Fatalf("EncodeMessage: %v", err)
	}

	var decoded Packet
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), &decoded); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}
	if decoded.IPVersion != uint32(2) {
		t.Fatalf("expected IPVersion 2, got %d", decoded.IPVersion)
	}
	if len(decoded.Samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(decoded.Samples))
	}

	sample, ok := decoded.Samples[0].(ExpandedFlowSample)
	if !ok {
		t.Fatal("expected ExpandedFlowSample")
	}
	if sample.Header.SampleSequenceNumber != uint32(100) {
		t.Fatalf("expected SampleSequenceNumber 100, got %d", sample.Header.SampleSequenceNumber)
	}
	if sample.Header.SourceIdValue != uint32(99) {
		t.Fatalf("expected SourceIdValue 99, got %d", sample.Header.SourceIdValue)
	}
	if sample.FlowRecordsCount != uint32(3) {
		t.Fatalf("expected FlowRecordsCount 3, got %d", sample.FlowRecordsCount)
	}
	if len(sample.Records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(sample.Records))
	}

	ipv4, ok := sample.Records[0].Data.(SampledIPv4)
	if !ok {
		t.Fatal("expected SampledIPv4")
	}
	if !bytes.Equal(ipv4.SrcIP, utils.IPAddress{192, 0, 2, 10}) {
		t.Fatalf("expected SrcIP {192,0,2,10}, got %v", ipv4.SrcIP)
	}
	if ipv4.DstPort != uint32(443) {
		t.Fatalf("expected DstPort 443, got %d", ipv4.DstPort)
	}

	sw, ok := sample.Records[1].Data.(ExtendedSwitch)
	if !ok {
		t.Fatal("expected ExtendedSwitch")
	}
	if sw.DstVlan != uint32(200) {
		t.Fatalf("expected DstVlan 200, got %d", sw.DstVlan)
	}

	gw, ok := sample.Records[2].Data.(ExtendedGateway)
	if !ok {
		t.Fatal("expected ExtendedGateway")
	}
	if !bytes.Equal(gw.NextHop, utils.IPAddress{203, 0, 113, 1}) {
		t.Fatalf("expected NextHop {203,0,113,1}, got %v", gw.NextHop)
	}
	if !reflect.DeepEqual(gw.ASPath, []uint32{64515, 64516}) {
		t.Fatalf("expected ASPath %v, got %v", []uint32{64515, 64516}, gw.ASPath)
	}
	if !reflect.DeepEqual(gw.Communities, []uint32{100, 200}) {
		t.Fatalf("expected Communities %v, got %v", []uint32{100, 200}, gw.Communities)
	}
}

func TestEncodeDecodeSFlowCounterSample(t *testing.T) {
	t.Parallel()
	packet := Packet{
		Version:        5,
		IPVersion:      1,
		AgentIP:        utils.IPAddress{198, 51, 100, 1},
		SubAgentId:     1,
		SequenceNumber: 2,
		Uptime:         3,
		Samples: []interface{}{
			CounterSample{
				Header: SampleHeader{
					Format:               SAMPLE_FORMAT_COUNTER,
					SampleSequenceNumber: 7,
					SourceIdType:         0,
					SourceIdValue:        8,
				},
				CounterRecordsCount: 2,
				Records: []CounterRecord{
					{
						Data: IfCounters{
							IfIndex:            1,
							IfType:             2,
							IfSpeed:            1000,
							IfDirection:        1,
							IfStatus:           3,
							IfInOctets:         100,
							IfInUcastPkts:      10,
							IfInMulticastPkts:  11,
							IfInBroadcastPkts:  12,
							IfInDiscards:       13,
							IfInErrors:         14,
							IfInUnknownProtos:  15,
							IfOutOctets:        200,
							IfOutUcastPkts:     20,
							IfOutMulticastPkts: 21,
							IfOutBroadcastPkts: 22,
							IfOutDiscards:      23,
							IfOutErrors:        24,
							IfPromiscuousMode:  1,
						},
					},
					{
						Data: EthernetCounters{
							Dot3StatsAlignmentErrors:           1,
							Dot3StatsFCSErrors:                 2,
							Dot3StatsSingleCollisionFrames:     3,
							Dot3StatsMultipleCollisionFrames:   4,
							Dot3StatsSQETestErrors:             5,
							Dot3StatsDeferredTransmissions:     6,
							Dot3StatsLateCollisions:            7,
							Dot3StatsExcessiveCollisions:       8,
							Dot3StatsInternalMacTransmitErrors: 9,
							Dot3StatsCarrierSenseErrors:        10,
							Dot3StatsFrameTooLongs:             11,
							Dot3StatsInternalMacReceiveErrors:  12,
							Dot3StatsSymbolErrors:              13,
						},
					},
				},
			},
		},
	}

	encoded, err := EncodeMessage(&packet)
	if err != nil {
		t.Fatalf("EncodeMessage: %v", err)
	}

	var decoded Packet
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), &decoded); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}
	if len(decoded.Samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(decoded.Samples))
	}

	sample, ok := decoded.Samples[0].(CounterSample)
	if !ok {
		t.Fatal("expected CounterSample")
	}
	if sample.CounterRecordsCount != uint32(2) {
		t.Fatalf("expected CounterRecordsCount 2, got %d", sample.CounterRecordsCount)
	}
	if len(sample.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(sample.Records))
	}

	ifc, ok := sample.Records[0].Data.(IfCounters)
	if !ok {
		t.Fatal("expected IfCounters")
	}
	if ifc.IfSpeed != uint64(1000) {
		t.Fatalf("expected IfSpeed 1000, got %d", ifc.IfSpeed)
	}
	if ifc.IfOutErrors != uint32(24) {
		t.Fatalf("expected IfOutErrors 24, got %d", ifc.IfOutErrors)
	}

	eth, ok := sample.Records[1].Data.(EthernetCounters)
	if !ok {
		t.Fatal("expected EthernetCounters")
	}
	if eth.Dot3StatsSymbolErrors != uint32(13) {
		t.Fatalf("expected Dot3StatsSymbolErrors 13, got %d", eth.Dot3StatsSymbolErrors)
	}
}

func TestEncodeDecodeSFlowDropSample(t *testing.T) {
	t.Parallel()
	packet := Packet{
		Version:        5,
		IPVersion:      1,
		AgentIP:        utils.IPAddress{192, 168, 119, 184},
		SubAgentId:     100000,
		SequenceNumber: 3,
		Uptime:         12350,
		Samples: []interface{}{
			DropSample{
				Header: SampleHeader{
					Format:               SAMPLE_FORMAT_DROP,
					SampleSequenceNumber: 2,
					SourceIdType:         0,
					SourceIdValue:        1,
				},
				Drops:            256,
				Input:            1,
				Output:           2,
				Reason:           1,
				FlowRecordsCount: 3,
				Records: []FlowRecord{
					{Data: EgressQueue{Queue: 42}},
					{Data: ExtendedACL{Number: 7, Name: "foo!", Direction: 2}},
					{Data: ExtendedFunction{Symbol: "dropper"}},
				},
			},
		},
	}

	encoded, err := EncodeMessage(&packet)
	if err != nil {
		t.Fatalf("EncodeMessage: %v", err)
	}

	var decoded Packet
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), &decoded); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}
	if len(decoded.Samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(decoded.Samples))
	}

	sample, ok := decoded.Samples[0].(DropSample)
	if !ok {
		t.Fatal("expected DropSample")
	}
	if sample.Reason != uint32(1) {
		t.Fatalf("expected Reason 1, got %d", sample.Reason)
	}
	if len(sample.Records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(sample.Records))
	}
	if sample.Records[0].Data != (EgressQueue{Queue: 42}) {
		t.Fatalf("expected EgressQueue{42}, got %v", sample.Records[0].Data)
	}
	if sample.Records[1].Data != (ExtendedACL{Number: 7, Name: "foo!", Direction: 2}) {
		t.Fatalf("expected ExtendedACL{7,foo!,2}, got %v", sample.Records[1].Data)
	}
	if sample.Records[2].Data != (ExtendedFunction{Symbol: "dropper"}) {
		t.Fatalf("expected ExtendedFunction{dropper}, got %v", sample.Records[2].Data)
	}
}
