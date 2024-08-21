package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/webitel/cases/internal/app"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
)

func Run() {
	log := wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  wlog.LevelDebug,
	})

	wlog.RedirectStdLog(log)
	wlog.InitGlobalLogger(log)

	config, appErr := loadConfig()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}

	application, appErr := app.New(config)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	initSignals(application)
	appErr = application.Start()
	wlog.Critical(appErr.Error())
}

func initSignals(application *app.App) {
	wlog.Info("initializing stop signals")
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
		wlog.Info("got kill signal, service gracefully stopped!")
		os.Exit(0)
	}
}

func loadConfig() (*model.AppConfig, model.AppError) {
	var appConfig model.AppConfig

	// Load from command-line flags
	dataSource := flag.String("data_source", "", "Data source")
	consul := flag.String("consul", "", "Host to consul")
	grpcAddr := flag.String("grpc_addr", "", "Public grpc address with port")
	consulID := flag.String("id", "", "Service id")
	flag.Parse()

	// Load from environment variables if flags are not provided
	if *dataSource == "" {
		*dataSource = os.Getenv("DATA_SOURCE")
	}
	if *consul == "" {
		*consul = os.Getenv("CONSUL")
	}
	if *grpcAddr == "" {
		*grpcAddr = os.Getenv("GRPC_ADDR")
	}
	if *consulID == "" {
		*consulID = os.Getenv("CONSUL_ID")
	}

	// Combined log statement for debugging
	wlog.Debug(fmt.Sprintf("Configuration - Data Source: %s, Consul: %s, gRPC Addr: %s, Consul ID: %s",
		*dataSource, *consul, *grpcAddr, *consulID))

	// Set the configuration struct fields
	appConfig.Database = &model.DatabaseConfig{
		Url: *dataSource,
	}
	appConfig.Consul = &model.ConsulConfig{
		Id:            *consulID,
		Address:       *consul,
		PublicAddress: *grpcAddr,
	}

	// Check if any required field is missing
	if appConfig.Database.Url == "" {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.missing_data_source", "Data source is required")
	}
	if appConfig.Consul.Id == "" {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.missing_id", "Service id is required")
	}
	if appConfig.Consul.Address == "" {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.missing_consul", "Consul address is required")
	}
	if appConfig.Consul.PublicAddress == "" {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.missing_grpc_addr", "gRPC address is required")
	}

	return &appConfig, nil
}
