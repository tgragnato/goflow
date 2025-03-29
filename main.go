package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	// decoders
	"github.com/tgragnato/goflow/decoders/netflow"
	"github.com/tgragnato/goflow/geoip"
	"github.com/tgragnato/goflow/sampler"

	// various formatters
	"github.com/tgragnato/goflow/format"
	_ "github.com/tgragnato/goflow/format/binary"
	_ "github.com/tgragnato/goflow/format/json"
	_ "github.com/tgragnato/goflow/format/text"

	// various transports
	"github.com/tgragnato/goflow/transport"
	_ "github.com/tgragnato/goflow/transport/file"
	_ "github.com/tgragnato/goflow/transport/syslog"

	// various producers
	"github.com/tgragnato/goflow/producer"
	protoproducer "github.com/tgragnato/goflow/producer/proto"
	rawproducer "github.com/tgragnato/goflow/producer/raw"

	// core libraries
	"github.com/tgragnato/goflow/metrics"
	"github.com/tgragnato/goflow/utils"
	"github.com/tgragnato/goflow/utils/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

var (
	ListenAddresses = flag.String("listen", "sflow://:6343,netflow://:2055", "listen addresses")

	Produce   = flag.String("produce", "sample", "Producer method (sample or raw)")
	Format    = flag.String("format", "json", fmt.Sprintf("Choose the format (available: %s)", strings.Join(format.GetFormats(), ", ")))
	Transport = flag.String("transport", "file", fmt.Sprintf("Choose the transport (available: %s)", strings.Join(transport.GetTransports(), ", ")))

	ErrCnt = flag.Int("err.cnt", 10, "Maximum errors per batch for muting")
	ErrInt = flag.Duration("err.int", time.Second*10, "Maximum errors interval for muting")

	Addr = flag.String("addr", ":8080", "HTTP server address")

	TemplatePath = flag.String("templates.path", "/templates", "NetFlow/IPFIX templates list")

	MappingFile = flag.String("mapping", "", "Configuration file for custom mappings")

	GeoipASN = flag.String("geoip.asn", "GeoLite2-ASN.mmdb", "Path to GeoIP ASN database")
	GeoipCC  = flag.String("geoip.cc", "GeoLite2-Country.mmdb", "Path to GeoIP Country database")
)

func LoadMapping(f io.Reader) (*protoproducer.ProducerConfig, error) {
	config := &protoproducer.ProducerConfig{}
	dec := yaml.NewDecoder(f)
	err := dec.Decode(config)
	return config, err
}

func main() {
	flag.Parse()
	geoip.Init(*GeoipASN, *GeoipCC)
	sampler.Init()

	formatter, err := format.FindFormat(*Format)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	transporter, err := transport.FindTransport(*Transport)
	if err != nil {
		fmt.Printf("error transporter: %s\n", err)
		os.Exit(1)
	}

	var flowProducer producer.ProducerInterface
	// instanciate a producer
	// unlike transport and format, the producer requires extensive configurations and can be chained
	switch *Produce {
	case "sample":
		var cfgProducer *protoproducer.ProducerConfig
		if *MappingFile != "" {
			f, err := os.Open(*MappingFile)
			if err != nil {
				fmt.Printf("error opening mapping: %s\n", err)
				os.Exit(1)
			}
			cfgProducer, err = LoadMapping(f)
			f.Close()
			if err != nil {
				fmt.Printf("error loading mapping: %s\n", err)
				os.Exit(1)
			}
		}

		cfgm, err := cfgProducer.Compile() // converts configuration into a format that can be used by a protobuf producer
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		flowProducer, err = protoproducer.CreateProtoProducer(cfgm, protoproducer.CreateSamplingSystem)
		if err != nil {
			slog.Error("error producer", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case "raw":
		flowProducer = &rawproducer.RawProducer{}
	default:
		fmt.Printf("producer %s does not exist\n", *Produce)
		os.Exit(1)
	}

	// intercept panic and generate an error
	flowProducer = debug.WrapPanicProducer(flowProducer)
	// wrap producer with Prometheus metrics
	flowProducer = metrics.WrapPromProducer(flowProducer)

	wg := &sync.WaitGroup{}

	var collecting bool
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/__health", func(wr http.ResponseWriter, r *http.Request) {
		if !collecting {
			wr.WriteHeader(http.StatusServiceUnavailable)
			if _, err := wr.Write([]byte("Not OK\n")); err != nil {
				fmt.Printf("error writing HTTP: %s\n", err)
			}
		} else {
			wr.WriteHeader(http.StatusOK)
			if _, err := wr.Write([]byte("OK\n")); err != nil {
				fmt.Printf("error writing HTTP: %s\n", err)
			}
		}
	})
	srv := http.Server{
		Addr:              *Addr,
		ReadHeaderTimeout: time.Second * 5,
	}
	if *Addr != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("HTTP server error: %s\n", err)
				os.Exit(1)
			}
		}()
	}

	fmt.Println("starting GoFlow")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var receivers []*utils.UDPReceiver
	var pipes []utils.FlowPipe

	q := make(chan bool)
	for _, listenAddress := range strings.Split(*ListenAddresses, ",") {
		listenAddrUrl, err := url.Parse(listenAddress)
		if err != nil {
			fmt.Printf("error parsing address: %s\n", err)
			os.Exit(1)
		}
		numSockets := 1
		if listenAddrUrl.Query().Has("count") {
			if numSocketsTmp, err := strconv.ParseUint(listenAddrUrl.Query().Get("count"), 10, 64); err != nil {
				fmt.Printf("error parsing count of sockets in URL: %s\n", err)
				os.Exit(1)
			} else if numSocketsTmp > math.MaxInt32 {
				numSockets = math.MaxInt32
			} else {
				numSockets = int(numSocketsTmp)
			}
		}
		if numSockets == 0 {
			numSockets = 1
		}

		var numWorkers int
		if listenAddrUrl.Query().Has("workers") {
			if numWorkersTmp, err := strconv.ParseUint(listenAddrUrl.Query().Get("workers"), 10, 64); err != nil {
				fmt.Printf("error parsing workers in URL: %s\n", err)
				os.Exit(1)
			} else if numWorkersTmp > math.MaxInt32 {
				numWorkers = math.MaxInt32
			} else {
				numWorkers = int(numWorkersTmp)
			}
		}
		if numWorkers == 0 && numSockets < math.MaxInt32/2 {
			numWorkers = numSockets * 2
		} else {
			numWorkers = math.MaxInt32
		}

		var isBlocking bool
		if listenAddrUrl.Query().Has("blocking") {
			if isBlocking, err = strconv.ParseBool(listenAddrUrl.Query().Get("blocking")); err != nil {
				fmt.Printf("error parsing blocking in URL: %s\n", err)
				os.Exit(1)
			}
		}

		var queueSize int
		if listenAddrUrl.Query().Has("queue_size") {
			if queueSizeTmp, err := strconv.ParseUint(listenAddrUrl.Query().Get("queue_size"), 10, 64); err != nil {
				fmt.Printf("error parsing queue_size in URL: %s\n", err)
				os.Exit(1)
			} else if queueSizeTmp > math.MaxInt32 {
				queueSize = math.MaxInt32
			} else {
				queueSize = int(queueSizeTmp)
			}
		} else if !isBlocking {
			queueSize = 1000000
		}

		hostname := listenAddrUrl.Hostname()
		port, err := strconv.ParseUint(listenAddrUrl.Port(), 10, 64)
		if err != nil {
			fmt.Printf("Port %s could not be converted to integer\n", listenAddrUrl.Port())
			return
		}

		logFields := map[string]interface{}{
			"scheme":     listenAddrUrl.Scheme,
			"hostname":   hostname,
			"port":       port,
			"count":      numSockets,
			"workers":    numWorkers,
			"blocking":   isBlocking,
			"queue_size": queueSize,
		}
		fmt.Printf("starting collection: \n%v\n", logFields)

		cfg := &utils.UDPReceiverConfig{
			Sockets:          numSockets,
			Workers:          numWorkers,
			QueueSize:        queueSize,
			Blocking:         isBlocking,
			ReceiverCallback: metrics.NewReceiverMetric(),
		}
		recv, err := utils.NewUDPReceiver(cfg)
		if err != nil {
			fmt.Printf("error creating UDP receiver: %s\n", err)
			os.Exit(1)
		}

		cfgPipe := &utils.PipeConfig{
			Format:           formatter,
			Transport:        transporter,
			Producer:         flowProducer,
			NetFlowTemplater: metrics.NewDefaultPromTemplateSystem, // wrap template system to get Prometheus info
		}

		var decodeFunc utils.DecoderFunc
		var p utils.FlowPipe
		switch scheme := listenAddrUrl.Scheme; scheme {
		case "sflow":
			p = utils.NewSFlowPipe(cfgPipe)
		case "netflow":
			p = utils.NewNetFlowPipe(cfgPipe)
		case "flow":
			p = utils.NewFlowPipe(cfgPipe)
		default:
			fmt.Printf("scheme %s does not exist\n", scheme)
			return
		}

		decodeFunc = p.DecodeFlow
		// intercept panic and generate error
		decodeFunc = debug.PanicDecoderWrapper(decodeFunc)
		// wrap decoder with Prometheus metrics
		decodeFunc = metrics.PromDecoderWrapper(decodeFunc, listenAddrUrl.Scheme)
		pipes = append(pipes, p)

		bm := utils.NewBatchMute(*ErrInt, *ErrCnt)

		if port < 1 || port > 65535 {
			fmt.Printf("Port %d is out of range\n", port)
			os.Exit(1)
		}

		// starts receivers
		// the function either returns an error
		if err := recv.Start(hostname, int(port), decodeFunc); err != nil {
			fmt.Printf("error starting: %s\n", err)
			os.Exit(1)
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case <-q:
						return
					case err := <-recv.Errors():
						if errors.Is(err, net.ErrClosed) {
							continue
						} else if !errors.Is(err, netflow.ErrorTemplateNotFound) && !errors.Is(err, debug.PanicError) {
							fmt.Printf("closed receiver: %s\n", err)
							continue
						}

						muted, skipped := bm.Increment()
						if muted && skipped == 0 {
							fmt.Println("too many receiver messages, muting")
						} else if !muted && skipped > 0 {
							fmt.Printf("skipped %d receiver messages\n", skipped)
						} else if !muted {
							if errors.Is(err, netflow.ErrorTemplateNotFound) {
								fmt.Printf("template error: %s\n", err)
							} else if errors.Is(err, debug.PanicError) {
								var pErrMsg *debug.PanicErrorMessage
								if errors.As(err, &pErrMsg) {
									fmt.Printf("intercepted panic (%s): \n\n%v\n", pErrMsg.Msg, pErrMsg.Stacktrace)
								} else {
									fmt.Printf("intercepted panic: %s\n", err)
								}
							}
						}

					}
				}
			}()
			receivers = append(receivers, recv)
		}
	}

	// special routine to handle kafka errors transmitted as a stream
	wg.Add(1)
	go func() {
		defer wg.Done()

		var transportErr <-chan error
		if transportErrorFct, ok := transporter.TransportDriver.(interface {
			Errors() <-chan error
		}); ok {
			transportErr = transportErrorFct.Errors()
		}

		bm := utils.NewBatchMute(*ErrInt, *ErrCnt)

		for {
			select {
			case <-q:
				return
			case err := <-transportErr:
				if err == nil {
					return
				}
				muted, skipped := bm.Increment()
				if muted && skipped == 0 {
					fmt.Println("too many transport errors, muting")
				} else if !muted && skipped > 0 {
					fmt.Printf("skipped %d transport errors\n", skipped)
				} else if !muted {
					fmt.Printf("transport error: %s\n", err)
				}

			}
		}
	}()

	collecting = true

	<-c

	collecting = false

	// stops receivers first, udp sockets will be down
	for _, recv := range receivers {
		if err := recv.Stop(); err != nil {
			fmt.Printf("error stopping receiver: %s\n", err)
		}
	}
	// then stop pipe
	for _, pipe := range pipes {
		pipe.Close()
	}
	// close producer
	flowProducer.Close()
	// close transporter (eg: flushes message to Kafka)
	if err := transporter.Close(); err != nil {
		fmt.Printf("closed transporter with error: %s\n", err)
	}

	// close http server (prometheus + health check)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting-down HTTP server: %s\n", err)
	}
	cancel()
	close(q) // close errors
	wg.Wait()
}
