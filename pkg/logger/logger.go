package logger

import (
	"io"
	"os"
	"path"
	"runtime/debug"

	"github.com/MihasBel/test-transactions-service/internal/app"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func errorStackMarshaller(_ error) interface{} {
	return string(debug.Stack())
}

// New creates a new instance for logger based on configuration
func New(cfg app.Configuration) zerolog.Logger {
	var writers []io.Writer

	if cfg.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}
	if cfg.FileLoggingEnabled {
		writers = append(writers, newRollingFile(&cfg))
	}

	mw := io.MultiWriter(writers...)
	zerolog.ErrorStackMarshaler = errorStackMarshaller
	l := zerolog.New(mw).With().Timestamp().Logger()

	l.Info().Msg("logging configured")

	return l
}
func newRollingFile(config *app.Configuration) io.Writer {
	if err := os.MkdirAll(config.LogDirectory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.LogDirectory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.LogDirectory, config.LogFilename),
		MaxBackups: config.LogMaxBackups, // files
		MaxSize:    config.LogMaxSize,    // megabytes
		MaxAge:     config.LogMaxAge,     // days
	}
}
