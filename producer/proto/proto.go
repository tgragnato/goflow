package protoproducer

import (
	"fmt"
	"sync"

	"github.com/tgragnato/goflow/decoders/netflow"
	"github.com/tgragnato/goflow/decoders/netflowlegacy"
	"github.com/tgragnato/goflow/decoders/sflow"
	"github.com/tgragnato/goflow/geoip"
	"github.com/tgragnato/goflow/producer"
	"github.com/tgragnato/goflow/sampler"
)

const LMS_TARGET_INDEX = "goflow.local"

type ProtoProducer struct {
	cfg                ProtoProducerConfig
	samplinglock       *sync.RWMutex
	sampling           map[string]SamplingRateSystem
	samplingRateSystem func() SamplingRateSystem
}

func (p *ProtoProducer) enrich(flowMessageSet []producer.ProducerMessage, cb func(msg *ProtoProducerMessage)) {
	for _, msg := range flowMessageSet {
		fmsg, ok := msg.(*ProtoProducerMessage)
		if !ok {
			continue
		}
		cb(fmsg)
	}
}

func (p *ProtoProducer) getSamplingRateSystem(args *producer.ProduceArgs) SamplingRateSystem {
	key := args.Src.Addr().String()
	p.samplinglock.RLock()
	sampling, ok := p.sampling[key]
	p.samplinglock.RUnlock()
	if !ok {
		sampling = p.samplingRateSystem()
		p.samplinglock.Lock()
		p.sampling[key] = sampling
		p.samplinglock.Unlock()
	}

	return sampling
}

func (p *ProtoProducer) Produce(msg interface{}, args *producer.ProduceArgs) (flowMessageSet []producer.ProducerMessage, err error) {
	tr := uint64(args.TimeReceived.UnixNano())
	sa, _ := args.SamplerAddress.Unmap().MarshalBinary()
	switch msgConv := msg.(type) {
	case *netflowlegacy.PacketNetFlowV5:
		flowMessageSet, err = ProcessMessageNetFlowLegacy(msgConv)

		p.enrich(flowMessageSet, func(fmsg *ProtoProducerMessage) {
			fmsg.TimeReceivedNs = tr
			fmsg.SamplerAddress = sa
			fmsg.LmsTargetIndex = LMS_TARGET_INDEX
			fmsg.SrcCountry = geoip.GetCountryByByteSlice(fmsg.SrcAddr)
			fmsg.SrcAsn, fmsg.SrcAs = geoip.GetASNByByteSlice(fmsg.SrcAddr)
			fmsg.DstCountry = geoip.GetCountryByByteSlice(fmsg.DstAddr)
			fmsg.DstAsn, fmsg.DstAs = geoip.GetASNByByteSlice(fmsg.DstAddr)
			fmsg.SamplerHostname = sampler.GetHostnameByByteSlice(fmsg.SamplerAddress)
			if len(fmsg.AsPath) == 0 {
				fmsg.AsPath = []uint32{fmsg.SrcAsn, 0, fmsg.DstAsn}
			}
		})
	case *netflow.NFv9Packet:
		samplingRateSystem := p.getSamplingRateSystem(args)
		flowMessageSet, err = ProcessMessageNetFlowV9Config(msgConv, samplingRateSystem, p.cfg)

		p.enrich(flowMessageSet, func(fmsg *ProtoProducerMessage) {
			fmsg.TimeReceivedNs = tr
			fmsg.SamplerAddress = sa
			fmsg.LmsTargetIndex = LMS_TARGET_INDEX
			fmsg.SrcCountry = geoip.GetCountryByByteSlice(fmsg.SrcAddr)
			fmsg.SrcAsn, fmsg.SrcAs = geoip.GetASNByByteSlice(fmsg.SrcAddr)
			fmsg.DstCountry = geoip.GetCountryByByteSlice(fmsg.DstAddr)
			fmsg.DstAsn, fmsg.DstAs = geoip.GetASNByByteSlice(fmsg.DstAddr)
			fmsg.SamplerHostname = sampler.GetHostnameByByteSlice(fmsg.SamplerAddress)
			if len(fmsg.AsPath) == 0 {
				fmsg.AsPath = []uint32{fmsg.SrcAsn, 0, fmsg.DstAsn}
			}
		})
	case *netflow.IPFIXPacket:
		samplingRateSystem := p.getSamplingRateSystem(args)
		flowMessageSet, err = ProcessMessageIPFIXConfig(msgConv, samplingRateSystem, p.cfg)

		p.enrich(flowMessageSet, func(fmsg *ProtoProducerMessage) {
			fmsg.TimeReceivedNs = tr
			fmsg.SamplerAddress = sa
			fmsg.LmsTargetIndex = LMS_TARGET_INDEX
			fmsg.SrcCountry = geoip.GetCountryByByteSlice(fmsg.SrcAddr)
			fmsg.SrcAsn, fmsg.SrcAs = geoip.GetASNByByteSlice(fmsg.SrcAddr)
			fmsg.DstCountry = geoip.GetCountryByByteSlice(fmsg.DstAddr)
			fmsg.DstAsn, fmsg.DstAs = geoip.GetASNByByteSlice(fmsg.DstAddr)
			fmsg.SamplerHostname = sampler.GetHostnameByByteSlice(fmsg.SamplerAddress)
			if len(fmsg.AsPath) == 0 {
				fmsg.AsPath = []uint32{fmsg.SrcAsn, 0, fmsg.DstAsn}
			}
		})
	case *sflow.Packet:
		flowMessageSet, err = ProcessMessageSFlowConfig(msgConv, p.cfg)

		p.enrich(flowMessageSet, func(fmsg *ProtoProducerMessage) {
			fmsg.TimeReceivedNs = tr
			fmsg.TimeFlowStartNs = tr
			fmsg.TimeFlowEndNs = tr
			fmsg.LmsTargetIndex = LMS_TARGET_INDEX
			fmsg.SrcCountry = geoip.GetCountryByByteSlice(fmsg.SrcAddr)
			fmsg.SrcAsn, fmsg.SrcAs = geoip.GetASNByByteSlice(fmsg.SrcAddr)
			fmsg.DstCountry = geoip.GetCountryByByteSlice(fmsg.DstAddr)
			fmsg.DstAsn, fmsg.DstAs = geoip.GetASNByByteSlice(fmsg.DstAddr)
			fmsg.SamplerHostname = sampler.GetHostnameByByteSlice(fmsg.SamplerAddress)
			if len(fmsg.AsPath) == 0 {
				fmsg.AsPath = []uint32{fmsg.SrcAsn, 0, fmsg.DstAsn}
			}
		})
	default:
		return flowMessageSet, fmt.Errorf("flow not recognized")
	}

	p.enrich(flowMessageSet, func(fmsg *ProtoProducerMessage) {
		fmsg.formatter = p.cfg.GetFormatter()
	})
	return flowMessageSet, err
}

func (p *ProtoProducer) Commit(flowMessageSet []producer.ProducerMessage) {
	for _, fmsg := range flowMessageSet {
		protoMessagePool.Put(fmsg)
	}
}

func (p *ProtoProducer) Close() {}

func CreateProtoProducer(cfg ProtoProducerConfig, samplingRateSystem func() SamplingRateSystem) (producer.ProducerInterface, error) {
	return &ProtoProducer{
		cfg:                cfg,
		samplinglock:       &sync.RWMutex{},
		sampling:           make(map[string]SamplingRateSystem),
		samplingRateSystem: samplingRateSystem,
	}, nil
}
