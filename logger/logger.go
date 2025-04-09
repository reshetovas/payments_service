package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(logLevel string) {
	zerolog.TimeFieldFormat = time.RFC3339 //time format RFC3339  (2025-03-31T12:00:00Z)

	switch strings.ToLower(logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	//setting up a global logger
	//zerolog.New - create new logger

	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	multi := io.MultiWriter(console, file)

	log.Logger = zerolog.New(multi).
		With().
		Timestamp().                       //add timestamp
		Str("service", "payment-service"). //add a general attribute
		Logger()
}

// package logger

// import (
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

// type LoggerStruct struct {
// 	logger zerolog.Logger
// }

// func NewLoggerStruct(logLevel string) *LoggerStruct {
// 	zerolog.TimeFieldFormat = time.RFC3339 //time format RFC3339  (2025-03-31T12:00:00Z)

// 	switch strings.ToLower(logLevel) {
// 	case "debug":
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	case "info":
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	case "warn":
// 		zerolog.SetGlobalLevel(zerolog.WarnLevel)
// 	case "error":
// 		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
// 	default:
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	}

// 	//setting up a global logger
// 	//zerolog.New - create new logger
// 	log.Logger = zerolog.New(os.Stdout).
// 		With().
// 		Timestamp().                       //add timestamp
// 		Str("service", "payment-service"). //add a general attribute
// 		Logger()

// 	return &LoggerStruct{
// 		logger: log.Logger,
// 	}
// }

// func (l *LoggerStruct) Fatal() *zerolog.Event {
// 	return l.logger.Fatal()
// }

// func (l *LoggerStruct) Error() *zerolog.Event {
// 	return l.logger.Error()
// }

// func (l *LoggerStruct) Warn() *zerolog.Event {
// 	return l.logger.Warn()
// }

// func (l *LoggerStruct) Info() *zerolog.Event {
// 	return l.logger.Info()
// }

// func (l *LoggerStruct) Debug() *zerolog.Event {
// 	return l.logger.Debug()
// }

// func (l *LoggerStruct) Trace() *zerolog.Event {
// 	return l.logger.Trace()
// }
