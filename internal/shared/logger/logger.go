package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Environment string
	LogLevel    string
	Format      string
}

func NewLogger(cfg Config) zerolog.Logger {

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer = os.Stdout
	if cfg.Format == "console" || cfg.Environment == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05.000",
		}
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Str("env", cfg.Environment).
		Logger()
}
