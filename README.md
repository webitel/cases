# Webitel Cases

## Available Flags and Environment Variables

### General
- `-config_file` (Flag) → `CASES_CONFIG_FILE` (Env)  
  _Configuration file in JSON format_ (default: "")

### Database
- `-data_source` (Flag) → `DATA_SOURCE` (Env)  
  _Data source connection string_ (default: "")

### Consul
- `-consul` (Flag) → `CONSUL` (Env)  
  _Host to Consul_ (default: "")
- `-grpc_addr` (Flag) → `GRPC_ADDR` (Env)  
  _Public gRPC address with port_ (default: "")
- `-id` (Flag) → `CONSUL_ID` (Env)  
  _Service ID_ (default: "")

### RabbitMQ
- `-amqp` (Flag) → `MICRO_BROKER_ADDRESS` (Env)  
  _AMQP connection URL_ (default: "")

### Trigger Watcher
- `-trigger_watcher_exchange` (Flag) → `TRIGGER_WATCHER_EXCHANGE_NAME` (Env)  
  _Watcher exchange name_ (default: "cases")
- `-trigger_watcher_topic` (Flag) → `TRIGGER_WATCHER_TOPIC_NAME` (Env)  
  _Watcher topic name_ (default: "*")
- `-trigger_watch_enabled` (Flag) → `TRIGGER_WATCHER_ENABLED` (Env)  
  _Enable Trigger Watcher_ (default: `true`)

### Logger Watcher
- `-logger_watch_enabled` (Flag) → `LOGGER_WATCHER_ENABLED` (Env)  
  _Enable Logger Watcher_ (default: `true`)

### FTS Watcher
- `-fts_watch_enabled` (Flag) → `FTS_WATCHER_ENABLED` (Env)  
  _Enable FTS Watcher_ (default: `false`)

### Global Watcher Control
- `-watchers_enabled` (Flag) → `WATCHERS_ENABLED` (Env)  
  _Enable all watchers (highest priority control)_ (default: `true`)
