package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	cerr "github.com/webitel/cases/internal/errors"
)

const defaultResolutionIntervalSec int64 = 5

// AppConfig and nested config structs...
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
	Url string `json:"url"`
}

type DatabaseConfig struct {
	Url string `json:"url"`
}

type TriggerWatcherConfig struct {
	ExchangeName            string `json:"exchange"`
	TopicName               string `json:"topic"`
	Enabled                 bool   `json:"enabled"`
	ResolutionCheckInterval int64  `json:"resolution_check_interval_sec"`
}

type FtsWatcherConfig struct {
	Enabled bool `json:"enabled"`
}

type LoggerWatcherConfig struct {
	Enabled bool `json:"enabled"`
}

type ConsulConfig struct {
	Id            string `json:"id"`
	Address       string `json:"address"`
	PublicAddress string `json:"publicAddress"`
}

func LoadConfig() (*AppConfig, error) {
	bindFlagsAndEnv()

	configFile := getConfigFilePath()
	if configFile != "" {
		if err := loadFromFile(configFile); err != nil {
			return nil, err
		}
	}

	cfg := buildAppConfig(configFile)
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindFlagsAndEnv() {
	pflag.String("config_file", "", "Configuration file in JSON format")
	pflag.String("data_source", "", "Data source")
	pflag.String("consul", "", "Host to consul")
	pflag.String("grpc_addr", "", "Public grpc address with port")
	pflag.String("id", "", "Service id")
	pflag.String("amqp", "", "AMQP connection URL")
	pflag.String("trigger_watcher_exchange", "cases", "Exchange name")
	pflag.String("trigger_watcher_topic", "*", "Queue name")
	pflag.Bool("trigger_watch_enabled", true, "Watcher enabled")
	pflag.Int64("resolution_check_interval_sec", defaultResolutionIntervalSec, "Interval between resolution checks")
	pflag.Bool("logger_watch_enabled", true, "Watcher enabled")
	pflag.Bool("fts_watch_enabled", false, "Watcher enabled")
	pflag.Bool("watchers_enabled", true, "Enable all watchers")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicit mapping
	_ = viper.BindEnv("id", "CONSUL_ID")
	_ = viper.BindEnv("amqp", "MICRO_BROKER_ADDRESS")
}

func getConfigFilePath() string {
	file := viper.GetString("config_file")
	if file == "" {
		file = os.Getenv("CASES_CONFIG_FILE")
	}

	return file
}

func loadFromFile(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		return cerr.NewInternalError(
			"cases.main.load_config",
			fmt.Sprintf("could not load config file: %s", err.Error()))
	}

	return nil
}

func buildAppConfig(file string) *AppConfig {
	return &AppConfig{
		File:     file,
		Database: &DatabaseConfig{Url: viper.GetString("data_source")},
		Consul: &ConsulConfig{
			Id:            viper.GetString("id"),
			Address:       viper.GetString("consul"),
			PublicAddress: viper.GetString("grpc_addr"),
		},
		Rabbit: &RabbitConfig{Url: viper.GetString("amqp")},
		TriggerWatcher: &TriggerWatcherConfig{
			ExchangeName:            viper.GetString("trigger_watcher_exchange"),
			TopicName:               viper.GetString("trigger_watcher_topic"),
			Enabled:                 viper.GetBool("trigger_watch_enabled"),
			ResolutionCheckInterval: viper.GetInt64("resolution_check_interval_sec"),
		},
		LoggerWatcher:   &LoggerWatcherConfig{Enabled: viper.GetBool("logger_watch_enabled")},
		FtsWatcher:      &FtsWatcherConfig{Enabled: viper.GetBool("fts_watch_enabled")},
		WatchersEnabled: viper.GetBool("watchers_enabled"),
	}
}

func validateConfig(cfg *AppConfig) error {
	if cfg.Database.Url == "" {
		return cerr.NewInternalError("cases.main.missing_data_source", "Data source is required")
	}
	if cfg.Consul.Id == "" {
		return cerr.NewInternalError("cases.main.missing_id", "Service id is required")
	}
	if cfg.Consul.Address == "" {
		return cerr.NewInternalError("cases.main.missing_consul", "Consul address is required")
	}
	if cfg.Consul.PublicAddress == "" {
		return cerr.NewInternalError("cases.main.missing_grpc_addr", "gRPC address is required")
	}
	if cfg.Rabbit.Url == "" {
		return cerr.NewInternalError("cases.main.missing_rabbit_url", "Rabbit URL is required")
	}

	return nil
}
