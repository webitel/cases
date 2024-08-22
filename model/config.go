package model

import (
	"flag"
	"os"
)

type AppConfig struct {
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty" `
	// Log      *LogSettings    `json:"log,omitempty"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source"`
}

type ConsulConfig struct {
	Id            string `json:"id" flag:"id|1| Service tag"`
	Address       string `json:"address" flag:"consul|| Host to consul"`
	PublicAddress string `json:"publicAddress" flag:"grpc_addr|| Public grpc address with port"`
}

func LoadConfig() (*AppConfig, AppError) {
	var appConfig AppConfig

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
	// if *logLevel == "" {
	// 	*logLevel = os.Getenv("LOG_LVL")
	// }
	// if !*logJson {
	// 	*logJson = os.Getenv("LOG_JSON") == "true"
	// }
	// if !*logOtel {
	// 	*logOtel = os.Getenv("LOG_OTEL") == "true"
	// }
	// if *logFile == "" {
	// 	*logFile = os.Getenv("LOG_FILE")
	// }

	// Set the configuration struct fields
	appConfig.Database = &DatabaseConfig{
		Url: *dataSource,
	}
	appConfig.Consul = &ConsulConfig{
		Id:            *consulID,
		Address:       *consul,
		PublicAddress: *grpcAddr,
	}

	// Check if any required field is missing
	if appConfig.Database.Url == "" {
		return nil, NewInternalError("cases.main.missing_data_source", "Data source is required")
	}
	if appConfig.Consul.Id == "" {
		return nil, NewInternalError("cases.main.missing_id", "Service id is required")
	}
	if appConfig.Consul.Address == "" {
		return nil, NewInternalError("cases.main.missing_consul", "Consul address is required")
	}
	if appConfig.Consul.PublicAddress == "" {
		return nil, NewInternalError("cases.main.missing_grpc_addr", "gRPC address is required")
	}

	return &appConfig, nil
}
