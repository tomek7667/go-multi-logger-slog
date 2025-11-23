package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/tomek7667/go-multi-logger-slog/internal/utils"
)

func SetLogLevel() {
	level := utils.Getenv("LOG_LEVEL", "DEBUG")
	logLevel := slog.LevelError.Level()
	switch strings.ToLower(level) {
	case "error":
	case "warn":
		logLevel = slog.LevelWarn.Level()
	case "info":
		logLevel = slog.LevelInfo.Level()
	default:
		logLevel = slog.LevelDebug.Level()
	}
	var logger *slog.Logger
	if os.Getenv("LOKI_ENDPOINT") == "" {
		slog.Info("Using tint logging as LOKI_ENDPOINT is not specified")
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      logLevel,
			AddSource:  true,
			TimeFormat: time.DateTime,
		}))
	} else {
		slog.Info("Using Loki logger")
		lokiLogger := slog.New(
			NewLokiHandler(
				os.Getenv("LOKI_ENDPOINT"),
				os.Getenv("LOKI_USERNAME"),
				os.Getenv("LOKI_PASSWORD"),
				os.Getenv("APP_NAME"),
				slog.Level(logLevel),
			),
		)
		jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.Level(logLevel),
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "time" {
					tmS := a.Value.String()
					tm, err := time.Parse(time.RFC3339Nano, tmS)
					if err != nil {
						return a
					}
					return slog.Attr{
						Key:   "time",
						Value: slog.AnyValue(tm.Format(time.DateTime)),
					}
				}
				return a
			},
		}))
		logger = slog.New(NewMultiHandler(lokiLogger.Handler(), jsonLogger.Handler()))
	}
	slog.SetDefault(logger)
	slog.Debug("Logger initialized", "level", logLevel)
}
