package cmd

import (
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/cases/internal/app"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath *string
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
	return
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
		wlog.Info(fmt.Sprintf("got kill signal, service gracefully stopped!"))
		os.Exit(0)
	}
}

func loadConfig() (*model.AppConfig, model.AppError) {
	var appConfig model.AppConfig

	configurator := configuration.New(
		&appConfig,
		// order of execution will be preserved:
		configuration.NewFlagProvider(),
		configuration.NewEnvProvider(),
		configuration.NewDefaultProvider(),
	)

	if err := configurator.InitValues(); err != nil {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.parse_fail", err.Error())
	}
	return &appConfig, nil
}
