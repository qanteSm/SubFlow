// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package pkg

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the global application logger
var (
	logger     zerolog.Logger
	loggerOnce sync.Once
)

// InitLogger initializes the global logger
func InitLogger(level string, pretty bool) {
	loggerOnce.Do(func() {
		// Parse log level
		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			lvl = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(lvl)

		// Create output writer
		var output = os.Stdout
		if pretty {
			output = os.Stdout
			logger = zerolog.New(zerolog.ConsoleWriter{
				Out:        output,
				TimeFormat: time.RFC3339,
			}).With().Timestamp().Caller().Logger()
		} else {
			logger = zerolog.New(output).With().Timestamp().Logger()
		}

		// Add architect signature to context
		logger = logger.With().Str("architect", "Muhammet-Ali-Buyuk").Logger()
	})
}

// GetLogger returns the global logger instance
func GetLogger() zerolog.Logger {
	return logger
}

// Debug logs a debug message
func Debug(msg string) {
	logger.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	logger.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	logger.Warn().Msg(msg)
}

// Error logs an error message
func Error(msg string, err error) {
	logger.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error) {
	logger.Fatal().Err(err).Msg(msg)
}

// WithFields returns a logger with additional fields
func WithFields(fields map[string]interface{}) zerolog.Logger {
	ctx := logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}
