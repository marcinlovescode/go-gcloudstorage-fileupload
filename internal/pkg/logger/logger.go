package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type ZerologLogger struct {
	logger *zerolog.Logger
}

func NewZerologLogger(level string) *ZerologLogger {
	var logLevel zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		logLevel = zerolog.ErrorLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &ZerologLogger{
		logger: &logger,
	}
}

func (logger *ZerologLogger) Debug(message interface{}, args ...interface{}) {
	logger.msg("debug", message, args...)
}

func (logger *ZerologLogger) Info(message string, args ...interface{}) {
	logger.log(message, args...)
}

func (logger *ZerologLogger) Warn(message string, args ...interface{}) {
	logger.log(message, args...)
}

func (logger *ZerologLogger) Error(message interface{}, args ...interface{}) {
	if logger.logger.GetLevel() == zerolog.DebugLevel {
		logger.Debug(message, args...)
	}

	logger.msg("error", message, args...)
}

func (logger *ZerologLogger) Fatal(message interface{}, args ...interface{}) {
	logger.msg("fatal", message, args...)

	os.Exit(1)
}

func (logger *ZerologLogger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.logger.Info().Msg(message)
	} else {
		logger.logger.Info().Msgf(message, args...)
	}
}

func (logger *ZerologLogger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		logger.log(msg.Error(), args...)
	case string:
		logger.log(msg, args...)
	default:
		logger.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
