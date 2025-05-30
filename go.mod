module github.com/webitel/cases

go 1.23.0

//replace buf.build/gen/go/webitel/cases/grpc/go/_gogrpc v1.5.1-20241105122241-a1d959115d35.1 => ./api/cases

replace github.com/webitel/proto/gen => github.com/webitel/custom/proto/gen v0.0.0-20250507094703-36254da8d7d5

require (
	buf.build/gen/go/webitel/webitel-go/grpc/go v1.5.1-20250121132426-2c80753dfee6.2
	buf.build/gen/go/webitel/webitel-go/protocolbuffers/go v1.36.1-20250121132426-2c80753dfee6.1
	github.com/gammazero/deque v1.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1
	github.com/jackc/pgtype v1.14.4
	github.com/webitel/custom v0.0.0-20250507094703-36254da8d7d5
	github.com/webitel/logger v0.0.0-20250131144320-fa79d885252a
	github.com/webitel/proto/gen v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250204164813-702378808489
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240920164238-5a7b106cbb87.2 // indirect
	buf.build/gen/go/webitel/logger/grpc/go v1.5.1-20250128105802-aaacc0377b27.2 // indirect
	buf.build/gen/go/webitel/logger/protocolbuffers/go v1.36.4-20250128105802-aaacc0377b27.1 // indirect
	cel.dev/expr v0.19.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/cel-go v0.22.1 // indirect
	github.com/grafana/otel-profiling-go v0.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/webitel/webitel-go-kit/logging/wlog v0.0.0-20240807083813-0853fbc06218 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
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
	go.opentelemetry.io/otel/sdk/metric v1.32.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/webitel/webitel-go-kit v0.0.21-0.20250522134942-6461799e3a09
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.11.0
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250204164813-702378808489
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/bufbuild/protovalidate-go v0.7.2
	github.com/hashicorp/consul/api v1.31.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v5 v5.7.2
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/mbobakov/grpc-consul-resolver v1.5.3
	github.com/nicksnyder/go-i18n v1.10.3
	github.com/rabbitmq/amqp091-go v1.10.0
	go.opentelemetry.io/contrib/bridges/otelslog v0.11.0
	go.opentelemetry.io/otel/log v0.12.2 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0
	google.golang.org/grpc v1.70.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
