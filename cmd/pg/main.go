package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/MihasBel/test-transactions-servise/adapters/pg"
	"github.com/MihasBel/test-transactions-servise/internal/app"
	"github.com/MihasBel/test-transactions-servise/pkg/logger"
	"github.com/google/uuid"
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
	db := pg.New(app.Config, l)
	if err := db.Start(context.Background()); err != nil {
		log.Error().Err(err)
	}
	/*err := db.PlaceTransaction(context.Background(), model.Transaction{
		ID:          uuid.New(),
		UserID:      uuid.MustParse("c3bb416e-9a47-11ed-a8fc-0242ac120002"),
		Amount:      -100,
		CreatedAt:   time.Now(),
		Status:      0,
		Description: "in processing",
	})*/
	tran, err := db.GetTransactionByID(context.Background(), uuid.MustParse("81f0a860-576d-4b7a-b56a-695442b065f9"))
	if err != nil {
		log.Error().Err(err)
	}
	log.Debug().Msg(fmt.Sprintln(tran))
	if err := db.Stop(context.TODO()); err != nil {
		log.Error().Err(err)
	}

}
