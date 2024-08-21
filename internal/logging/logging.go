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

// Setup initializes the logging configuration based on the provided LogSettings.
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
	if logConfig.Json {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	logger := slog.New(handler)

	// If OTEL is enabled, setup OTEL logging
	if logConfig.Otel {
		ctx := context.Background()
		shutdown, err := otelsdk.Setup(
			ctx,
			otelsdk.WithResource(service),
			otelsdk.WithLogLevel(log.SeverityDebug),
		)
		defer shutdown(ctx)

		if err != nil {
			logger.ErrorContext(ctx, "OTel setup failed", slog.String("error", err.Error()))
			os.Exit(1)
		}

		logger.InfoContext(ctx, "OTel setup successful", slog.Bool("success", true))
	}

	// Set the global logger
	slog.SetDefault(logger)
}
