package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/webitel/cases/model"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Setup(service *resource.Resource) func(context.Context) error {
	// Initialize the context and OTel setup
	ctx := context.Background()
	shutdown, err := otelsdk.Setup(
		ctx,
		otelsdk.WithResource(service),
		otelsdk.WithLogLevel(log.SeverityDebug),
	)

	// Initialize the OTel logger
	stdlog := otelslog.NewLogger(model.APP_SERVICE_NAME)

	// Check if OTel setup failed
	if err != nil {
		stdlog.ErrorContext(
			ctx, "OTel setup",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Set slog default to stdlog [OTel] logger
	slog.SetDefault(stdlog)

	// OTEL setup successful
	slog.InfoContext(ctx, "OTel setup successful")

	// Return the shutdown function
	return shutdown
}
