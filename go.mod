module github.com/webitel/cases

go 1.23.0

//replace (
//	buf.build/gen/go/webitel/cases/grpc/go/_gogrpc v1.5.1-20241105122241-a1d959115d35.1 => ./api/cases
//)

require (
	buf.build/gen/go/webitel/cases/grpc/go v1.5.1-20241106124357-9fdd668c1b5c.1
	buf.build/gen/go/webitel/cases/protocolbuffers/go v1.35.1-20241105122241-a1d959115d35.1
	buf.build/gen/go/webitel/general/protocolbuffers/go v1.35.1-20240527133232-b30251c3af1c.1
	buf.build/gen/go/webitel/webitel-go/protocolbuffers/go v1.34.2-20240725111140-8ade8003e454.2
	github.com/webitel/cases/api/cases v0.0.0-20241105162143-aa9f2dbc1af7
	google.golang.org/protobuf v1.35.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240920164238-5a7b106cbb87.2 // indirect
	buf.build/gen/go/grpc-ecosystem/grpc-gateway/protocolbuffers/go v1.35.1-20240617172850-a48fcebcf8f1.1 // indirect
	buf.build/gen/go/webitel/general/grpc/go v1.5.1-20240527133232-b30251c3af1c.1 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/cel-go v0.21.0 // indirect
	github.com/grafana/otel-profiling-go v0.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/webitel/webitel-go-kit/logging/wlog v0.0.0-20240807083813-0853fbc06218 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.28.0 // indirect
	go.opentelemetry.io/contrib/samplers/jaegerremote v0.22.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.0.0-20240805233418-127d068751eb // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.4.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.5.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241104194629-dd2ea8efbc28 // indirect
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/google/uuid v1.6.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/webitel/webitel-go-kit v0.0.21-0.20241105101307-954c5cba127c
	go.opentelemetry.io/otel v1.29.0
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	buf.build/gen/go/webitel/webitel-go/grpc/go v1.5.1-20240725111140-8ade8003e454.1
	github.com/Masterminds/squirrel v1.5.4
	github.com/bufbuild/protovalidate-go v0.7.2
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/hashicorp/consul/api v1.29.2
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v5 v5.6.0
	github.com/lib/pq v1.10.9
	github.com/mbobakov/grpc-consul-resolver v1.5.3
	github.com/nicksnyder/go-i18n v1.10.3
	github.com/rabbitmq/amqp091-go v1.10.0
	go.opentelemetry.io/contrib/bridges/otelslog v0.3.0
	go.opentelemetry.io/otel/log v0.5.0
	go.opentelemetry.io/otel/sdk v1.29.0
	google.golang.org/grpc v1.67.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
