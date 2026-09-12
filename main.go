package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tgragnato.it/goflow/geoip"
	"tgragnato.it/goflow/pkg/goflow2/app"
	"tgragnato.it/goflow/pkg/goflow2/config"
	"tgragnato.it/goflow/sampler"

	_ "tgragnato.it/goflow/format/binary"
	_ "tgragnato.it/goflow/format/json"
	_ "tgragnato.it/goflow/format/text"
	_ "tgragnato.it/goflow/transport/file"
	_ "tgragnato.it/goflow/transport/syslog"
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
