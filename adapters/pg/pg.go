package pg

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MihasBel/test-transactions-servise/internal/app"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PG represents connection to PostgreSQL DB
type PG struct {
	log  *zerolog.Logger
	gorm *gorm.DB
	dsn  string
	cfg  app.Configuration
}

// New create new PG
func New(cfg app.Configuration, l zerolog.Logger) *PG {

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	return &PG{
		log: &l,
		dsn: dsn,
		cfg: cfg,
	}
}

// Start starts PG adapter
func (pg *PG) Start(context.Context) error {

	gormConfig := &gorm.Config{}

	if pg.cfg.DebugSQL {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: pg.dsn,
	}), gormConfig)
	if err != nil {
		pg.log.Error().Err(err).Msg("failed to open GORM")
		return err
	}

	pg.gorm = gormDB

	pg.log.Info().
		Str("DBName", pg.cfg.DBName).
		Str("Host", pg.cfg.Host).
		Int("Port", pg.cfg.Port).
		Str("User", pg.cfg.User).
		Msg("db connected")

	return nil
}

// Stop stops PG adapter
func (pg *PG) Stop(_ context.Context) error {
	sql, err := pg.gorm.DB()
	if err != nil {
		pg.log.Error().Err(err).Msg("failed to get db")
	}
	if err := sql.Close(); err != nil {
		pg.log.Error().Err(err).Msg("failed to close db")
	}

	return nil
}
