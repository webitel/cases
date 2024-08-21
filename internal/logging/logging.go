package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/webitel/cases/model"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Setup(logConfig *model.LogSettings, service *resource.Resource) {
	var logLevel slog.Level
	switch logConfig.Lvl {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Set up the logging handler
	var handler slog.Handler

	if logConfig.File == "" {
		handler = defaultHandler(logConfig, logLevel)
	} else {
		file, err := os.OpenFile(logConfig.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			handler = defaultHandler(logConfig, logLevel)
		} else {
			if logConfig.Json {
				handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: logLevel})
			} else {
				handler = slog.NewTextHandler(file, &slog.HandlerOptions{Level: logLevel})
			}
		}
	}

	slog.SetDefault(slog.New(handler))

	// OTEL setup
	if logConfig.Otel {
		ctx := context.Background()
		shutdown, err := otelsdk.Setup(
			ctx,
			otelsdk.WithResource(service),
			otelsdk.WithLogLevel(log.SeverityDebug),
		)
		defer shutdown(ctx)

		if err != nil {
			slog.Error("cases.logging.setup", slog.String("error", "OTel setup failed"))
			os.Exit(1)
		}

		slog.Info("cases.logging.setup", slog.String("message", "OTel setup successful"))
	}

	slog.Info("cases.logging.setup", slog.String("message", "Logging setup complete"))
}

// defaultHandler creates a default slog handler based on the configuration.
func defaultHandler(logConfig *model.LogSettings, logLevel slog.Level) slog.Handler {
	if logConfig.Json {
		slog.Info("cases.logging.setup", slog.String("message", "Logging in JSON format"))
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	}
	slog.Info("cases.logging.setup", slog.String("message", "Logging in text format"))
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
}
