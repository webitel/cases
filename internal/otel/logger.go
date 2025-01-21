package logging

import (
	"context"
	"log/slog"
	"os"

	slogutil "github.com/webitel/webitel-go-kit/otel/log/bridge/slog"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/resource"

	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/stdout"
)

// Setup initializes OpenTelemetry with slog logging and returns a shutdown function
func Setup(service *resource.Resource) func(context.Context) error {
	// Retrieve log level from the environment, default to info
	var verbose slog.LevelVar
	verbose.Set(slog.LevelInfo)
	if input := os.Getenv("OTEL_LOG_LEVEL"); input != "" {
		_ = verbose.UnmarshalText([]byte(input))
	}

	// Initialize the context and OpenTelemetry setup
	ctx := context.Background()
	shutdown, err := otelsdk.Configure(
		ctx,
		otelsdk.WithResource(service),
		otelsdk.WithLogBridge(func() {
			// Redirect slog.Default() to OpenTelemetry
			stdlog := slog.New(
				slogutil.WithLevel(
					&verbose,                    // Filter level for otelslog.Handler
					otelslog.NewHandler("slog"), // otelslog Handler for OpenTelemetry
				),
			)
			slog.SetDefault(stdlog)
		}),
	)

	// Initialize the logger
	log := slog.Default()

	// Error handling if OpenTelemetry setup fails
	if err != nil {
		log.ErrorContext(
			ctx, "OpenTelemetry setup failed",
			"error", err,
		)
		os.Exit(1)
	}

	// Log setup success message
	log.InfoContext(ctx, "OpenTelemetry setup successful")

	return shutdown
}
