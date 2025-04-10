package model

import (
	"encoding/json"
	"flag"
	"fmt"
	cerr "github.com/webitel/cases/internal/errors"
	"os"
	"strconv"
)

type AppConfig struct {
	File            string                `json:"-"`
	Rabbit          *RabbitConfig         `json:"rabbit,omitempty"`
	Database        *DatabaseConfig       `json:"database,omitempty"`
	Consul          *ConsulConfig         `json:"consul,omitempty"`
	TriggerWatcher  *TriggerWatcherConfig `json:"trigger_watcher,omitempty"`
	FtsWatcher      *FtsWatcherConfig     `json:"fts_watcher,omitempty"`
	LoggerWatcher   *LoggerWatcherConfig  `json:"logger_watcher,omitempty"`
	WatchersEnabled bool                  `json:"watchers_enabled,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url" flag:"amqp|| AMQP connection"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source"`
}

type TriggerWatcherConfig struct {
	ExchangeName            string `json:"exchange" flag:"trigger_watcher_exchange || watcher exchange"`
	TopicName               string `json:"topic" flag:"trigger_watcher_topic || watcher topic"`
	Enabled                 bool   `json:"enabled" flag:"trigger_watch_enabled || watch_enabled"`
	ResolutionCheckInterval int64  `json:"resolution_check_interval_sec" flag:"resolution_check_interval_sec || watch_enabled"`
}

type FtsWatcherConfig struct {
	Enabled bool `json:"enabled" flag:"fts_watch_enabled || watch_enabled"`
}

type LoggerWatcherConfig struct {
	Enabled bool `json:"enabled" flag:"logger_watch_enabled || watch_enabled"`
}

type ConsulConfig struct {
	Id            string `json:"id" flag:"id|1| Service tag"`
	Address       string `json:"address" flag:"consul|| Host to consul"`
	PublicAddress string `json:"publicAddress" flag:"grpc_addr|| Public grpc address with port"`
}

func LoadConfig() (*AppConfig, error) { // Change to return standard error
	var appConfig AppConfig

	// TODO :: refactor processing default values
	// Load from command-line flags
	dataSource := flag.String("data_source", "", "Data source")
	consul := flag.String("consul", "", "Host to consul")
	grpcAddr := flag.String("grpc_addr", "", "Public grpc address with port")
	consulID := flag.String("id", "", "Service id")
	rabbitURL := flag.String("amqp", "", "AMQP connection URL")
	triggerConfig := new(TriggerWatcherConfig)
	flag.StringVar(&triggerConfig.ExchangeName, "trigger_watcher_exchange", "", "Exchange name")
	flag.StringVar(&triggerConfig.TopicName, "trigger_watcher_topic", "", "Queue name")
	flag.BoolVar(&triggerConfig.Enabled, "trigger_watch_enabled", true, "Watcher enabled")
	flag.Int64Var(&triggerConfig.ResolutionCheckInterval, "resolution_check_interval_sec", 5, "The period, measured in seconds, between consecutive checks for resolution updates")

	loggerConfig := new(LoggerWatcherConfig)
	flag.BoolVar(&loggerConfig.Enabled, "logger_watch_enabled", true, "Watcher enabled")

	ftsConfig := new(FtsWatcherConfig)
	flag.BoolVar(&ftsConfig.Enabled, "fts_watch_enabled", false, "Watcher enabled")

	flag.BoolVar(&appConfig.WatchersEnabled, "watchers_enabled", true, "Flag controls all watchers and has the highest priority")

	// add possibility to load config from file
	flag.StringVar(&appConfig.File, "config_file", "", "Configuration file in JSON format")

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
	if *rabbitURL == "" {
		*rabbitURL = os.Getenv("MICRO_BROKER_ADDRESS")
	}

	if triggerConfig.ExchangeName == "" {
		value := "cases"
		if env := os.Getenv("TRIGGER_WATCHER_EXCHANGE_NAME"); env != "" {
			value = env
		}
		triggerConfig.ExchangeName = value
	}

	if triggerConfig.TopicName == "" {
		value := "*"
		if env := os.Getenv("TRIGGER_WATCHER_TOPIC_NAME"); env != "" {
			value = env
		}
		triggerConfig.TopicName = value
	}

	if env := os.Getenv("TRIGGER_WATCHER_ENABLED"); env != "" {
		triggerConfig.Enabled = env == "1" || env == "true"
	}

	if env := os.Getenv("TRIGGER_RESOLUTION_CHECK_INTERVAL_SEC"); env != "" {
		i, _ := strconv.ParseInt(env, 10, 64)
		if i > 0 {
			triggerConfig.ResolutionCheckInterval = i
		}
	}

	if env := os.Getenv("LOGGER_WATCHER_ENABLED"); env != "" {
		loggerConfig.Enabled = env == "1" || env == "true"
	}

	if env := os.Getenv("FTS_WATCHER_ENABLED"); env != "" {
		ftsConfig.Enabled = env == "1" || env == "true"
	}

	if env := os.Getenv("WATCHERS_ENABLED"); env != "" {
		appConfig.WatchersEnabled = env == "1" || env == "true"
	}

	// Set the configuration struct fields
	appConfig.Database = &DatabaseConfig{
		Url: *dataSource,
	}
	appConfig.Consul = &ConsulConfig{
		Id:            *consulID,
		Address:       *consul,
		PublicAddress: *grpcAddr,
	}
	appConfig.Rabbit = &RabbitConfig{
		Url: *rabbitURL,
	}
	appConfig.TriggerWatcher = triggerConfig
	appConfig.LoggerWatcher = loggerConfig
	appConfig.FtsWatcher = ftsConfig

	// trying to load config from file
	if appConfig.File == "" {
		appConfig.File = os.Getenv("CASES_CONFIG_FILE")
	}

	if appConfig.File != "" {
		configData, err := os.ReadFile(appConfig.File)
		if err != nil {
			return nil, cerr.NewInternalError("cases.main.load_config", fmt.Sprintf("could not load config file: %s", err.Error()))
		}
		err = json.Unmarshal(configData, &appConfig)
		if err != nil {
			return nil, cerr.NewInternalError("cases.main.parse_config", fmt.Sprintf("could not parse config file: %s", err.Error()))
		}
	}

	// Check if any required field is missing
	if appConfig.Database.Url == "" {
		return nil, cerr.NewInternalError("cases.main.missing_data_source", "Data source is required")
	}
	if appConfig.Consul.Id == "" {
		return nil, cerr.NewInternalError("cases.main.missing_id", "Service id is required")
	}
	if appConfig.Consul.Address == "" {
		return nil, cerr.NewInternalError("cases.main.missing_consul", "Consul address is required")
	}
	if appConfig.Consul.PublicAddress == "" {
		return nil, cerr.NewInternalError("cases.main.missing_grpc_addr", "gRPC address is required")
	}
	if appConfig.Rabbit.Url == "" {
		return nil, cerr.NewInternalError("cases.main.missing_rabbit_url", "Rabbit URL is required")
	}

	return &appConfig, nil
}
