package model

import (
	"flag"
	"os"

	conferr "github.com/webitel/cases/internal/error"
)

type AppConfig struct {
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url" flag:"amqp|| AMQP connection"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source"`
}

type ConsulConfig struct {
	Id            string `json:"id" flag:"id|1| Service tag"`
	Address       string `json:"address" flag:"consul|| Host to consul"`
	PublicAddress string `json:"publicAddress" flag:"grpc_addr|| Public grpc address with port"`
}

func LoadConfig() (*AppConfig, error) { // Change to return standard error
	var appConfig AppConfig

	// Load from command-line flags
	dataSource := flag.String("data_source", "", "Data source")
	consul := flag.String("consul", "", "Host to consul")
	grpcAddr := flag.String("grpc_addr", "", "Public grpc address with port")
	consulID := flag.String("id", "", "Service id")
	rabbitURL := flag.String("amqp", "", "AMQP connection URL")

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

	// Check if any required field is missing
	if appConfig.Database.Url == "" {
		return nil, conferr.NewConfigError("cases.main.missing_data_source", "Data source is required")
	}
	if appConfig.Consul.Id == "" {
		return nil, conferr.NewConfigError("cases.main.missing_id", "Service id is required")
	}
	if appConfig.Consul.Address == "" {
		return nil, conferr.NewConfigError("cases.main.missing_consul", "Consul address is required")
	}
	if appConfig.Consul.PublicAddress == "" {
		return nil, conferr.NewConfigError("cases.main.missing_grpc_addr", "gRPC address is required")
	}
	if appConfig.Rabbit.Url == "" {
		return nil, conferr.NewConfigError("cases.main.missing_rabbit_url", "Rabbit URL is required")
	}

	return &appConfig, nil
}
