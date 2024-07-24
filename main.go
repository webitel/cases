package main

import (
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/cases/app"
	"github.com/webitel/cases/model"
	"github.com/webitel/wlog"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath *string
)

func main() {
	//flagDefine()
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
	// * Create an application layer
	app, appErr := app.New(config)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	initSignals(app)
	appErr = app.Start()
	wlog.Critical(appErr.Error())
	return

}

func initSignals(app *app.App) {
	wlog.Info("initializing stop signals")
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go func() {
		for {
			s := <-sigchnl
			handleSignals(s, app)
		}
	}()

}

func handleSignals(signal os.Signal, app *app.App) {
	if signal == syscall.SIGTERM || signal == syscall.SIGINT || signal == syscall.SIGKILL {
		app.Stop()
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
