package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MihasBel/test-transactions-service/internal/app"
	"github.com/MihasBel/test-transactions-service/pkg/logger"
	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "configs/local/env.json", "Config file path")
	flag.Parse()

	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&app.Config, configPath); err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	l := logger.New(app.Config)
	log.Logger = l

	startCtx, startCancel := context.WithTimeout(context.Background(), time.Duration(app.Config.StartTimeout)*time.Second)
	defer startCancel()
	a := app.New(app.Config)
	if err := a.Start(startCtx); err != nil {
		log.Fatal().Err(err).Msg("start error")
	}
	log.Info().Msg("application started")

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Duration(app.Config.StartTimeout)*time.Second)
	defer stopCancel()
	if err := a.Stop(stopCtx); err != nil {
		log.Fatal().Err(err).Msg("stop error")
	}
	log.Info().Msg("service is down")
}
