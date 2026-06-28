package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tgragnato/goflow/geoip"
	"github.com/tgragnato/goflow/pkg/goflow2/app"
	"github.com/tgragnato/goflow/pkg/goflow2/config"
	"github.com/tgragnato/goflow/sampler"

	_ "github.com/tgragnato/goflow/format/binary"
	_ "github.com/tgragnato/goflow/format/json"
	_ "github.com/tgragnato/goflow/format/text"
	_ "github.com/tgragnato/goflow/transport/file"
	_ "github.com/tgragnato/goflow/transport/syslog"
)

func main() {
	cfg := config.BindFlags(flag.CommandLine)
	flag.Parse()

	geoip.Init(cfg.GeoipASN, cfg.GeoipCC)
	sampler.Init()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		slog.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
