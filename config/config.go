package model

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	cerr "github.com/webitel/cases/internal/errors"
)

type AppConfig struct {
	File     string          `json:"-"`
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty"`
	Watcher  *WatcherConfig  `json:"watcher,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url" flag:"amqp|| AMQP connection"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source"`
}

type WatcherConfig struct {
	ExchangeName      string `json:"exchange" flag:"watcher_exchange || watcher exchange"`
	QueueName         string `json:"queue" flag:"watcher_queue || watcher queue"`
	TopicName         string `json:"topic" flag:"watcher_topic || watcher topic"`
	AMQPUser          string `json:"amqp_user" flag:"amqp_user || AMQP user"`
	QueuesMessagesTTL int    `json:"queues_messages_ttl" flag:"watcher_messages_ttl || Watcher queues messages TTL in milliseconds"`
	Enabled           bool   `json:"enabled" flag:"watch_enabled || watch_enabled"`
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
	watcher := new(WatcherConfig)
	flag.StringVar(&watcher.ExchangeName, "watcher_exchange", "", "Exchange name")
	flag.StringVar(&watcher.QueueName, "watcher_queue", "", "Queue name")
	flag.StringVar(&watcher.TopicName, "watcher_topic", "", "Queue name")
	flag.StringVar(&watcher.AMQPUser, "amqp_user", "", "AMQP user for publishing messages")
	flag.IntVar(&watcher.QueuesMessagesTTL, "watcher_messages_ttl", 0, "Watcher queues messages TTL in milliseconds")
	flag.BoolVar(&watcher.Enabled, "watch_enabled", true, "Watcher enabled")

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

	if watcher.ExchangeName == "" {
		value := "watcher_exchange"
		if env := os.Getenv("WATCHER_EXCHANGE_NAME"); env != "" {
			value = env
		}
		watcher.ExchangeName = value
	}

	if watcher.QueueName == "" {
		value := "watcher_queue"
		if env := os.Getenv("WATCHER_QUEUE_NAME"); env != "" {
			value = env
		}
		watcher.QueueName = value
	}

	if watcher.TopicName == "" {
		value := "trigger_case.*"
		if env := os.Getenv("WATCHER_TOPIC_NAME"); env != "" {
			value = env
		}
		watcher.TopicName = value
	}

	if env := os.Getenv("WATCHER_ENABLED"); env != "" {
		watcher.Enabled = env == "1" || env == "true"
	}

	if watcher.AMQPUser == "" {
		value := "webitel"
		if env := os.Getenv("WATCHER_AMQP_USER"); env != "" {
			value = env
		}
		watcher.AMQPUser = value
	}

	if watcher.QueuesMessagesTTL == 0 {
		value := 10000
		if env := os.Getenv("WATCHER_MESSAGES_TTL"); env != "" {
			value, _ = strconv.Atoi(env)
		}
		watcher.QueuesMessagesTTL = value
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
	appConfig.Watcher = watcher

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
