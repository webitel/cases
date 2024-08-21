package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/webitel/cases/internal/app"
	"github.com/webitel/cases/internal/logging"
	"github.com/webitel/cases/model"

	// ----- LOGGING -----//

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	// -------------------- plugin(s) -------------------- //
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/stdout"
)

func Run() {
	// Load configuration
	config, appErr := model.LoadConfig()
	if appErr != nil {
		slog.Error("cases.main.configuration_error", slog.String("error", appErr.Error()))
		return
	}

	// Initialize the application
	application, appErr := app.New(config)
	if appErr != nil {
		slog.Error("cases.main.application_initialization_error", slog.String("error", appErr.Error()))
		return
	}

	// Initialize signal handling for graceful shutdown
	initSignals(application)

	// ----- LOGGING -----//

	// slog + OTEL logging
	service := resource.NewSchemaless(
		semconv.ServiceName(model.APP_SERVICE_NAME),
		semconv.ServiceVersion(model.CurrentVersion),
		semconv.ServiceInstanceID(config.Consul.Id),
		semconv.ServiceNamespace(model.NAMESPACE_NAME),
	)
	logging.Setup(config.Log, service) // Use the logging package for setup

	// Log the configuration
	slog.Debug("cases.main.configuration_loaded",
		slog.String("data_source", config.Database.Url),
		slog.String("consul", config.Consul.Address),
		slog.String("grpc_address", config.Consul.PublicAddress),
		slog.String("consul_id", config.Consul.Id),
		slog.String("log_level", config.Log.Lvl),
		slog.Bool("log_json", config.Log.Json),
		slog.Bool("log_otel", config.Log.Otel),
		slog.String("log_file", config.Log.File),
	)

	// Start the application
	slog.Info("cases.main.starting_application")
	appErr = application.Start()
	if appErr != nil {
		slog.Error("cases.main.application_start_error", slog.String("error", appErr.Error()))
	} else {
		slog.Info("cases.main.application_started_successfully")
	}
}

func initSignals(application *app.App) {
	slog.Info("cases.main.initializing_stop_signals", slog.String("main", "initializing_stop_signals"))
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go func() {
		for {
			s := <-sigchnl
			handleSignals(s, application)
		}
	}()
}

func handleSignals(signal os.Signal, application *app.App) {
	if signal == syscall.SIGTERM || signal == syscall.SIGINT || signal == syscall.SIGKILL {
		application.Stop()
		slog.Info("cases.main.received_kill_signal", slog.String("signal", signal.String()), slog.String("status", "service gracefully stopped"))
		os.Exit(0)
	}
}
