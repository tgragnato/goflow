package netflow

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncodeDecodeNetFlowV9(t *testing.T) {
	t.Parallel()
	packet := NFv9Packet{
		Version:        9,
		Count:          4,
		SystemUptime:   100,
		UnixSeconds:    200,
		SequenceNumber: 300,
		SourceId:       400,
		FlowSets: []interface{}{
			TemplateFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 0},
				Records: []TemplateRecord{
					{
						TemplateId: 256,
						FieldCount: 2,
						Fields: []Field{
							{Type: NFV9_FIELD_IN_BYTES, Length: 4},
							{Type: NFV9_FIELD_IPV4_SRC_ADDR, Length: 4},
						},
					},
				},
			},
			NFv9OptionsTemplateFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 1},
				Records: []NFv9OptionsTemplateRecord{
					{
						TemplateId:   257,
						ScopeLength:  4,
						OptionLength: 4,
						Scopes: []Field{
							{Type: 1, Length: 4},
						},
						Options: []Field{
							{Type: NFV9_FIELD_SAMPLING_INTERVAL, Length: 4},
						},
					},
				},
			},
			DataFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 256},
				Records: []DataRecord{
					{
						Values: []DataField{
							{Type: NFV9_FIELD_IN_BYTES, Value: []byte{0, 0, 0, 10}},
							{Type: NFV9_FIELD_IPV4_SRC_ADDR, Value: []byte{192, 0, 2, 1}},
						},
					},
				},
			},
			OptionsDataFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 257},
				Records: []OptionsDataRecord{
					{
						ScopesValues: []DataField{
							{Type: 1, Value: []byte{0, 0, 0, 1}},
						},
						OptionsValues: []DataField{
							{Type: NFV9_FIELD_SAMPLING_INTERVAL, Value: []byte{0, 0, 3, 232}},
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

	store := newTestTemplateStore()
	ctx := FlowContext{RouterKey: "test-router"}
	var decoded NFv9Packet
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), store, ctx, &decoded, nil); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}

	if decoded.Version != packet.Version {
		t.Fatalf("expected Version %d, got %d", packet.Version, decoded.Version)
	}
	if decoded.Count != packet.Count {
		t.Fatalf("expected Count %d, got %d", packet.Count, decoded.Count)
	}
	if decoded.SystemUptime != packet.SystemUptime {
		t.Fatalf("expected SystemUptime %d, got %d", packet.SystemUptime, decoded.SystemUptime)
	}
	if decoded.UnixSeconds != packet.UnixSeconds {
		t.Fatalf("expected UnixSeconds %d, got %d", packet.UnixSeconds, decoded.UnixSeconds)
	}
	if decoded.SequenceNumber != packet.SequenceNumber {
		t.Fatalf("expected SequenceNumber %d, got %d", packet.SequenceNumber, decoded.SequenceNumber)
	}
	if decoded.SourceId != packet.SourceId {
		t.Fatalf("expected SourceId %d, got %d", packet.SourceId, decoded.SourceId)
	}
	if len(decoded.FlowSets) != 4 {
		t.Fatalf("expected 4 FlowSets, got %d", len(decoded.FlowSets))
	}

	templateSet, ok := decoded.FlowSets[0].(TemplateFlowSet)
	if !ok {
		t.Fatal("expected TemplateFlowSet")
	}
	if templateSet.Id != uint16(0) {
		t.Fatalf("expected Id 0, got %d", templateSet.Id)
	}
	wantFields := []Field{
		{Type: NFV9_FIELD_IN_BYTES, Length: 4},
		{Type: NFV9_FIELD_IPV4_SRC_ADDR, Length: 4},
	}
	if !reflect.DeepEqual(templateSet.Records[0].Fields, wantFields) {
		t.Fatalf("expected Fields %v, got %v", wantFields, templateSet.Records[0].Fields)
	}

	optionsTemplateSet, ok := decoded.FlowSets[1].(NFv9OptionsTemplateFlowSet)
	if !ok {
		t.Fatal("expected NFv9OptionsTemplateFlowSet")
	}
	if optionsTemplateSet.Id != uint16(1) {
		t.Fatalf("expected Id 1, got %d", optionsTemplateSet.Id)
	}
	if !reflect.DeepEqual(optionsTemplateSet.Records[0].Scopes, []Field{{Type: 1, Length: 4}}) {
		t.Fatalf("unexpected Scopes: %v", optionsTemplateSet.Records[0].Scopes)
	}
	if !reflect.DeepEqual(optionsTemplateSet.Records[0].Options, []Field{{Type: NFV9_FIELD_SAMPLING_INTERVAL, Length: 4}}) {
		t.Fatalf("unexpected Options: %v", optionsTemplateSet.Records[0].Options)
	}

	dataSet, ok := decoded.FlowSets[2].(DataFlowSet)
	if !ok {
		t.Fatal("expected DataFlowSet")
	}
	if dataSet.Id != uint16(256) {
		t.Fatalf("expected Id 256, got %d", dataSet.Id)
	}
	if !bytes.Equal(dataSet.Records[0].Values[0].Value.([]byte), []byte{0, 0, 0, 10}) {
		t.Fatalf("unexpected Values[0]: %v", dataSet.Records[0].Values[0].Value)
	}
	if !bytes.Equal(dataSet.Records[0].Values[1].Value.([]byte), []byte{192, 0, 2, 1}) {
		t.Fatalf("unexpected Values[1]: %v", dataSet.Records[0].Values[1].Value)
	}

	optionsDataSet, ok := decoded.FlowSets[3].(OptionsDataFlowSet)
	if !ok {
		t.Fatal("expected OptionsDataFlowSet")
	}
	if optionsDataSet.Id != uint16(257) {
		t.Fatalf("expected Id 257, got %d", optionsDataSet.Id)
	}
	if !bytes.Equal(optionsDataSet.Records[0].ScopesValues[0].Value.([]byte), []byte{0, 0, 0, 1}) {
		t.Fatalf("unexpected ScopesValues[0]: %v", optionsDataSet.Records[0].ScopesValues[0].Value)
	}
	if !bytes.Equal(optionsDataSet.Records[0].OptionsValues[0].Value.([]byte), []byte{0, 0, 3, 232}) {
		t.Fatalf("unexpected OptionsValues[0]: %v", optionsDataSet.Records[0].OptionsValues[0].Value)
	}
}

func TestEncodeDecodeIPFIX(t *testing.T) {
	t.Parallel()
	packet := IPFIXPacket{
		Version:             10,
		ExportTime:          123,
		SequenceNumber:      456,
		ObservationDomainId: 789,
		FlowSets: []interface{}{
			TemplateFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 2},
				Records: []TemplateRecord{
					{
						TemplateId: 300,
						FieldCount: 2,
						Fields: []Field{
							{Type: IPFIX_FIELD_octetDeltaCount, Length: 8},
							{PenProvided: true, Type: 4000, Length: 2, Pen: 32473},
						},
					},
				},
			},
			IPFIXOptionsTemplateFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 3},
				Records: []IPFIXOptionsTemplateRecord{
					{
						TemplateId:      301,
						FieldCount:      2,
						ScopeFieldCount: 1,
						Scopes: []Field{
							{Type: IPFIX_FIELD_observationDomainId, Length: 4},
						},
						Options: []Field{
							{Type: IPFIX_FIELD_samplingInterval, Length: 0xffff},
						},
					},
				},
			},
			DataFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 300},
				Records: []DataRecord{
					{
						Values: []DataField{
							{Type: IPFIX_FIELD_octetDeltaCount, Value: []byte{0, 0, 0, 0, 0, 0, 0, 5}},
							{PenProvided: true, Type: 4000, Pen: 32473, Value: []byte{0x12, 0x34}},
						},
					},
				},
			},
			OptionsDataFlowSet{
				FlowSetHeader: FlowSetHeader{Id: 301},
				Records: []OptionsDataRecord{
					{
						ScopesValues: []DataField{
							{Type: IPFIX_FIELD_observationDomainId, Value: []byte{0, 0, 3, 21}},
						},
						OptionsValues: []DataField{
							{Type: IPFIX_FIELD_samplingInterval, Value: []byte("abcde")},
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

	store := newTestTemplateStore()
	ctx := FlowContext{RouterKey: "test-router"}
	var decoded IPFIXPacket
	if err := DecodeMessageVersion(bytes.NewBuffer(encoded), store, ctx, nil, &decoded); err != nil {
		t.Fatalf("DecodeMessageVersion: %v", err)
	}

	if decoded.Version != uint16(10) {
		t.Fatalf("expected Version 10, got %d", decoded.Version)
	}
	if decoded.Length != uint16(len(encoded)) {
		t.Fatalf("expected Length %d, got %d", len(encoded), decoded.Length)
	}
	if decoded.ExportTime != packet.ExportTime {
		t.Fatalf("expected ExportTime %d, got %d", packet.ExportTime, decoded.ExportTime)
	}
	if decoded.SequenceNumber != packet.SequenceNumber {
		t.Fatalf("expected SequenceNumber %d, got %d", packet.SequenceNumber, decoded.SequenceNumber)
	}
	if decoded.ObservationDomainId != packet.ObservationDomainId {
		t.Fatalf("expected ObservationDomainId %d, got %d", packet.ObservationDomainId, decoded.ObservationDomainId)
	}
	if len(decoded.FlowSets) != 4 {
		t.Fatalf("expected 4 FlowSets, got %d", len(decoded.FlowSets))
	}

	templateSet, ok := decoded.FlowSets[0].(TemplateFlowSet)
	if !ok {
		t.Fatal("expected TemplateFlowSet")
	}
	if templateSet.Id != uint16(2) {
		t.Fatalf("expected Id 2, got %d", templateSet.Id)
	}
	wantFields := []Field{
		{Type: IPFIX_FIELD_octetDeltaCount, Length: 8},
		{PenProvided: true, Type: 4000, Length: 2, Pen: 32473},
	}
	if !reflect.DeepEqual(templateSet.Records[0].Fields, wantFields) {
		t.Fatalf("expected Fields %v, got %v", wantFields, templateSet.Records[0].Fields)
	}

	optionsTemplateSet, ok := decoded.FlowSets[1].(IPFIXOptionsTemplateFlowSet)
	if !ok {
		t.Fatal("expected IPFIXOptionsTemplateFlowSet")
	}
	if optionsTemplateSet.Id != uint16(3) {
		t.Fatalf("expected Id 3, got %d", optionsTemplateSet.Id)
	}
	if !reflect.DeepEqual(optionsTemplateSet.Records[0].Scopes, []Field{{Type: IPFIX_FIELD_observationDomainId, Length: 4}}) {
		t.Fatalf("unexpected Scopes: %v", optionsTemplateSet.Records[0].Scopes)
	}
	if !reflect.DeepEqual(optionsTemplateSet.Records[0].Options, []Field{{Type: IPFIX_FIELD_samplingInterval, Length: 0xffff}}) {
		t.Fatalf("unexpected Options: %v", optionsTemplateSet.Records[0].Options)
	}

	dataSet, ok := decoded.FlowSets[2].(DataFlowSet)
	if !ok {
		t.Fatal("expected DataFlowSet")
	}
	if !bytes.Equal(dataSet.Records[0].Values[0].Value.([]byte), []byte{0, 0, 0, 0, 0, 0, 0, 5}) {
		t.Fatalf("unexpected Values[0]: %v", dataSet.Records[0].Values[0].Value)
	}
	if !bytes.Equal(dataSet.Records[0].Values[1].Value.([]byte), []byte{0x12, 0x34}) {
		t.Fatalf("unexpected Values[1]: %v", dataSet.Records[0].Values[1].Value)
	}

	optionsDataSet, ok := decoded.FlowSets[3].(OptionsDataFlowSet)
	if !ok {
		t.Fatal("expected OptionsDataFlowSet")
	}
	if !bytes.Equal(optionsDataSet.Records[0].ScopesValues[0].Value.([]byte), []byte{0, 0, 3, 21}) {
		t.Fatalf("unexpected ScopesValues[0]: %v", optionsDataSet.Records[0].ScopesValues[0].Value)
	}
	if !bytes.Equal(optionsDataSet.Records[0].OptionsValues[0].Value.([]byte), []byte("abcde")) {
		t.Fatalf("unexpected OptionsValues[0]: %v", optionsDataSet.Records[0].OptionsValues[0].Value)
	}
}
