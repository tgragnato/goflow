package protoproducer

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

func TestProcessEthernet(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"86dd" // etype

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseEthernet(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseEthernet: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if flowMessage.Etype != uint32(0x86dd) {
		t.Fatalf("expected Etype 0x86dd, got 0x%x", flowMessage.Etype)
	}
}

func TestProcessDot1Q(t *testing.T) {
	t.Parallel()

	dataStr := "00140800"

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = Parse8021Q(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("Parse8021Q: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if flowMessage.VlanId != uint32(20) {
		t.Fatalf("expected VlanId 20, got %d", flowMessage.VlanId)
	}
	if flowMessage.Etype != uint32(0x0800) {
		t.Fatalf("expected Etype 0x0800, got 0x%x", flowMessage.Etype)
	}
}

func TestProcessMPLS(t *testing.T) {
	t.Parallel()

	dataStr := "000120ff" + // label 1
		"000101ff" // label 2

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseMPLS(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseMPLS: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if !reflect.DeepEqual(flowMessage.MplsLabel, []uint32{18, 16}) {
		t.Fatalf("expected MplsLabel %v, got %v", []uint32{18, 16}, flowMessage.MplsLabel)
	}
	if !reflect.DeepEqual(flowMessage.MplsTtl, []uint32{255, 255}) {
		t.Fatalf("expected MplsTtl %v, got %v", []uint32{255, 255}, flowMessage.MplsTtl)
	}
	//assert.Equal(t, uint32(0x800), flowMessage.Etype) // tested with next byte in whole packet
}

func TestProcessIPv4(t *testing.T) {
	t.Parallel()

	dataStr := "45000064" +
		"abab" + // id
		"0000ff01" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" // dst

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseIPv4(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseIPv4: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if !bytes.Equal(flowMessage.SrcAddr, []byte{10, 0, 0, 1}) {
		t.Fatalf("expected SrcAddr %v, got %v", []byte{10, 0, 0, 1}, flowMessage.SrcAddr)
	}
	if !bytes.Equal(flowMessage.DstAddr, []byte{10, 0, 0, 2}) {
		t.Fatalf("expected DstAddr %v, got %v", []byte{10, 0, 0, 2}, flowMessage.DstAddr)
	}
	if flowMessage.FragmentId != uint32(0xabab) {
		t.Fatalf("expected FragmentId 0xabab, got 0x%x", flowMessage.FragmentId)
	}
	if flowMessage.IpTtl != uint32(0xff) {
		t.Fatalf("expected IpTtl 0xff, got 0x%x", flowMessage.IpTtl)
	}
	if flowMessage.Proto != uint32(1) {
		t.Fatalf("expected Proto 1, got %d", flowMessage.Proto)
	}
}

func TestProcessIPv6(t *testing.T) {
	t.Parallel()

	dataStr := "6001010104d83a40" + // ipv6
		"fd010000000000000000000000000001" + // src
		"fd010000000000000000000000000002" // dst

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseIPv6(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseIPv6: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	wantSrc := []byte{0xfd, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	wantDst := []byte{0xfd, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}
	if !bytes.Equal(flowMessage.SrcAddr, wantSrc) {
		t.Fatalf("expected SrcAddr %v, got %v", wantSrc, flowMessage.SrcAddr)
	}
	if !bytes.Equal(flowMessage.DstAddr, wantDst) {
		t.Fatalf("expected DstAddr %v, got %v", wantDst, flowMessage.DstAddr)
	}
	if flowMessage.IpTtl != uint32(0x40) {
		t.Fatalf("expected IpTtl 0x40, got 0x%x", flowMessage.IpTtl)
	}
	if flowMessage.Proto != uint32(0x3a) {
		t.Fatalf("expected Proto 0x3a, got 0x%x", flowMessage.Proto)
	}
	if flowMessage.Ipv6FlowLabel != uint32(0x010101) {
		t.Fatalf("expected Ipv6FlowLabel 0x010101, got 0x%x", flowMessage.Ipv6FlowLabel)
	}
}

func TestProcessIPv6HeaderFragment(t *testing.T) {
	t.Parallel()

	dataStr := "3a000001a7882ea9"

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseIPv6HeaderFragment(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseIPv6HeaderFragment: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if flowMessage.FragmentId != uint32(2810719913) {
		t.Fatalf("expected FragmentId 2810719913, got %d", flowMessage.FragmentId)
	}
	if flowMessage.FragmentOffset != uint32(0) {
		t.Fatalf("expected FragmentOffset 0, got %d", flowMessage.FragmentOffset)
	}
}

func TestProcessIPv6HeaderRouting(t *testing.T) {
	t.Parallel()

	dataStr := "29060401020300102001baba0002e00200000000000000002001baba0001000000000000000000002001baba0003e0070000000000000000"

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseIPv6HeaderRouting(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseIPv6HeaderRouting: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))
}

func TestProcessICMP(t *testing.T) {
	t.Parallel()

	dataStr := "01018cf7000627c4"

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseICMP(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseICMP: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if flowMessage.IcmpType != uint32(1) {
		t.Fatalf("expected IcmpType 1, got %d", flowMessage.IcmpType)
	}
	if flowMessage.IcmpCode != uint32(1) {
		t.Fatalf("expected IcmpCode 1, got %d", flowMessage.IcmpCode)
	}
}

func TestProcessICMPv6(t *testing.T) {
	t.Parallel()

	dataStr := "8080f96508a4"

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	_, err = ParseICMPv6(&flowMessage, data, ParseConfig{})
	if err != nil {
		t.Fatalf("ParseICMPv6: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	if flowMessage.IcmpType != uint32(128) {
		t.Fatalf("expected IcmpType 128, got %d", flowMessage.IcmpType)
	}
	if flowMessage.IcmpCode != uint32(128) {
		t.Fatalf("expected IcmpCode 128, got %d", flowMessage.IcmpCode)
	}
}

func TestProcessPacketBase(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"8100" + // etype
		"00008847" + // 8021q
		"000120ff" + // mpls label 1
		"000101ff" + // mpls label 2
		"6000000004d83a40" + // ipv6
		"fd010000000000000000000000000001" + // src
		"fd010000000000000000000000000002" + // dst
		"8000f96508a4" // icmpv6

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	err = ParsePacket(&flowMessage, data, nil, nil)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	layers := []uint32{0, 6, 5, 2, 8}
	if len(flowMessage.LayerStack) != len(layers) {
		t.Fatalf("expected %d layers, got %d", len(layers), len(flowMessage.LayerStack))
	}
	for i, layer := range layers {
		if uint32(flowMessage.LayerStack[i]) != layer {
			t.Fatalf("layer[%d]: expected %d, got %d", i, layer, flowMessage.LayerStack[i])
		}
	}

	if flowMessage.Etype != uint32(0x86dd) {
		t.Fatalf("expected Etype 0x86dd, got 0x%x", flowMessage.Etype)
	}
}

func TestProcessPacketGRE(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"86dd" + // etype

		"6000000004d82f40" + // ipv6
		"fd010000000000000000000000000001" + // src
		"fd010000000000000000000000000002" + // dst

		"00000800" + // gre

		"45000064" + // ipv4
		"abab" + // id
		"0000ff01" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		"01018cf7000627c4" // icmp

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	err = ParsePacket(&flowMessage, data, nil, nil)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	layers := []uint32{0, 2, 9, 1, 7}
	if len(flowMessage.LayerStack) != len(layers) {
		t.Fatalf("expected %d layers, got %d", len(layers), len(flowMessage.LayerStack))
	}
	for i, layer := range layers {
		if uint32(flowMessage.LayerStack[i]) != layer {
			t.Fatalf("layer[%d]: expected %d, got %d", i, layer, flowMessage.LayerStack[i])
		}
	}

	if flowMessage.Etype != uint32(0x86dd) {
		t.Fatalf("expected Etype 0x86dd, got 0x%x", flowMessage.Etype)
	}
	if flowMessage.Proto != uint32(47) {
		t.Fatalf("expected Proto 47, got %d", flowMessage.Proto)
	}
	// todo: check addresses
}

type testProtoProducerMessage struct {
	ProtoProducerMessage
	t *testing.T
}

func (m *testProtoProducerMessage) MapCustom(key string, v []byte, cfg MappableField) error {
	m.t.Log("mapping", key, v)
	mc := MapConfigBase{
		Endianness: BigEndian,
		ProtoIndex: 999,
		ProtoType:  ProtoVarint,
		ProtoArray: false,
	}
	return m.ProtoProducerMessage.MapCustom(key, v, &mc)
}

func TestProcessPacketMapping(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"0800" + // etype

		"45000064" + // ipv4
		"abab" + // id
		"0000ff11" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		// udp
		"ff00" + // src port
		"0035" + // dst port
		"0010" + // length
		"ffff" + // csum

		"0000000000000000" // payload

	config := SFlowProducerConfig{
		Mapping: []SFlowMapField{
			SFlowMapField{
				Layer:  "udp",
				Offset: 48,
				Length: 16,

				Destination: "csum",
			},
		},
	}
	configm := mapFieldsSFlow(config.Mapping)

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	flowMessage := testProtoProducerMessage{
		t: t,
	}

	err = ParsePacket(&flowMessage, data, configm, nil)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))
}

func TestProcessPacketMappingEncap(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"86dd" + // etype

		"6001010104d82b40" + // ipv6
		"fd010000000000000000000000000001" + // src
		"fd010000000000000000000000000002" + // dst

		"04060401020300102001baba0002e00200000000000000002001baba0001000000000000000000002001baba0003e0070000000000000000" + // srv6

		"45000064" + // ipv4
		"abab" + // id
		"0000ff11" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		// udp
		"ff00" + // src port
		"0035" + // dst port
		"0010" + // length
		"ffff" + // csum

		"0000000000000000" // payload

	config := ProducerConfig{
		Formatter: FormatterConfig{
			Render: map[string]RendererID{
				"src_ip_encap": RendererIP,
				"dst_ip_encap": RendererIP,
			},
			Fields: []string{
				"src_ip_encap",
				"dst_ip_encap",
			},
			Protobuf: []ProtobufFormatterConfig{
				ProtobufFormatterConfig{
					Name:  "src_ip_encap",
					Index: 998,
					Type:  "string",
					Array: true,
				},
				ProtobufFormatterConfig{
					Name:  "dst_ip_encap",
					Index: 999,
					Type:  "string",
					Array: true,
				},
			},
		},
		SFlow: SFlowProducerConfig{
			Mapping: []SFlowMapField{
				SFlowMapField{
					Layer:        "ipv6",
					Offset:       64,
					Length:       128,
					Encapsulated: true,

					Destination: "src_ip_encap",
				},
				SFlowMapField{
					Layer:        "ipv6",
					Offset:       192,
					Length:       128,
					Encapsulated: true,

					Destination: "dst_ip_encap",
				},

				SFlowMapField{
					Layer:        "ipv4",
					Offset:       96,
					Length:       32,
					Encapsulated: true,

					Destination: "src_ip_encap",
				},
				SFlowMapField{
					Layer:        "ipv4",
					Offset:       128,
					Length:       32,
					Encapsulated: true,

					Destination: "dst_ip_encap",
				},
			},
		},
	}
	configm, _ := config.Compile()

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage
	flowMessage.formatter = configm.GetFormatter()

	err = configm.GetPacketMapper().ParsePacket(&flowMessage, data)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}

	b, _ := json.Marshal(&flowMessage.FlowMessage)
	t.Log(string(b))

	flowMessage.skipDelimiter = true
	b, _ = flowMessage.MarshalBinary()
	t.Log(base64.StdEncoding.EncodeToString(b))

	b, _ = flowMessage.MarshalJSON()
	t.Log(string(b))
}

func TestProcessPacketMappingPort(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"0800" + // etype

		"45000064" + // ipv4
		"abab" + // id
		"0000ff11" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		// udp
		"ff00" + // src port
		"0035" + // dst port
		"0015" + // length
		"ffff" + // csum

		"02a901000001000000000000146578616d706c6503636f6d0000010001" // dns packet

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage

	var domain []byte

	pe := NewBaseParserEnvironment()

	if err := pe.RegisterPort("udp", PortDirDst, 53, ParserInfo{
		Parser: func(flowMessage *ProtoProducerMessage, data []byte, pc ParseConfig) (res ParseResult, err error) {
			domain = data[13 : 13+11]
			flowMessage.AddLayer("Custom")
			t.Log("read DNS packet", string(domain))
			res.Size = len(data)
			return res, err
		},
	}); err != nil {
		t.Fatalf("RegisterPort: %v", err)
	}

	err = ParsePacket(&flowMessage, data, nil, pe)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}
	wantDomain := []byte{0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x03, 0x63, 0x6F, 0x6D}
	if !bytes.Equal(domain, wantDomain) {
		t.Fatalf("expected domain %v, got %v", wantDomain, domain)
	}
	if len(flowMessage.LayerSize) != 4 {
		t.Fatalf("expected 4 LayerSize entries, got %d", len(flowMessage.LayerSize))
	}
	if flowMessage.LayerSize[3] != uint32(29) {
		t.Fatalf("expected LayerSize[3] 29, got %d", flowMessage.LayerSize[3])
	}
}

func TestProcessPacketMappingGeneve(t *testing.T) {
	t.Parallel()

	dataStr := "005300000001" + // src mac
		"005300000002" + // dst mac
		"0800" + // etype

		"45000064" + // ipv4
		"abab" + // id
		"0000ff11" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		// udp
		"ff00" + // src port
		"17c1" + // dst port
		"0015" + // length
		"ffff" + // csum

		"0240655800000a00000080010000000c" + // geneve

		"005300000001" + // src mac
		"005300000002" + // dst mac
		"0800" + // etype

		"45000064" + // ipv4
		"abab" + // id
		"0000ff01" + // flag, ttl, proto
		"aaaa" + // csum
		"0a000001" + // src
		"0a000002" + // dst

		"01018cf7000627c4" // icmp

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("hex.DecodeString: %v", err)
	}

	var flowMessage ProtoProducerMessage

	pe := NewBaseParserEnvironment()

	gp, ok := pe.GetParser("geneve")
	if !ok {
		t.Fatal("GetParser(geneve) returned false")
	}

	if err := pe.RegisterPort("udp", PortDirBoth, 6081, gp); err != nil {
		t.Fatalf("RegisterPort: %v", err)
	}

	err = ParsePacket(&flowMessage, data, nil, pe)
	if err != nil {
		t.Fatalf("ParsePacket: %v", err)
	}

	layers := []uint32{0, 1, 4, 12, 0, 1, 7}
	if len(flowMessage.LayerStack) != len(layers) {
		t.Fatalf("expected %d layers, got %d", len(layers), len(flowMessage.LayerStack))
	}
	for i, layer := range layers {
		if uint32(flowMessage.LayerStack[i]) != layer {
			t.Fatalf("layer[%d]: expected %d, got %d", i, layer, flowMessage.LayerStack[i])
		}
	}
}
