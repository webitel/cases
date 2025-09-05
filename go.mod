module github.com/webitel/cases

go 1.24.0

toolchain go1.24.6

//replace buf.build/gen/go/webitel/cases/grpc/go/_gogrpc v1.5.1-20241105122241-a1d959115d35.1 => ./api/cases

replace github.com/webitel/proto/gen => github.com/webitel/custom/proto/gen v0.0.0-20250507094703-36254da8d7d5

require (
	buf.build/gen/go/webitel/webitel-go/grpc/go v1.5.1-20250121132426-2c80753dfee6.2
	buf.build/gen/go/webitel/webitel-go/protocolbuffers/go v1.36.1-20250121132426-2c80753dfee6.1
	github.com/gammazero/deque v1.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/jackc/pgtype v1.14.4
	github.com/webitel/custom v0.0.0-20250507094703-36254da8d7d5
	github.com/webitel/proto/gen v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250826171959-ef028d996bc1
	google.golang.org/protobuf v1.36.8
)

require (
	github.com/georgysavva/scany/v2 v2.1.4
	github.com/google/cel-go v0.26.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
	github.com/webitel/webitel-go-kit/cmd/protoc-gen-go-webitel v0.0.0-20240829153325-0ae7f6059b52
	github.com/webitel/webitel-go-kit/infra/fts_client v0.0.0-20250702140655-98a92d815fcb
	github.com/webitel/webitel-go-kit/infra/logger_client v0.0.0-20250702140655-98a92d815fcb
	github.com/webitel/webitel-go-kit/infra/otel v0.0.0-20250625090308-5d99e087fa32
	github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq v0.0.0-20250702140655-98a92d815fcb
	github.com/webitel/webitel-go-kit/pkg/errors v0.0.0-20250625095313-01c817ef50a5
	github.com/webitel/webitel-go-kit/pkg/etag v0.0.0-20250625095313-01c817ef50a5
	github.com/webitel/webitel-go-kit/pkg/filters v0.0.0-20250905143306-9f3c3d8c69c7
	github.com/webitel/webitel-go-kit/pkg/watcher v0.0.0-20250625090308-5d99e087fa32
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240920164238-5a7b106cbb87.2 // indirect
	cel.dev/expr v0.24.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.12.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.12.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.12.2 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
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
	go.opentelemetry.io/otel v1.36.0
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/bufbuild/protovalidate-go v0.7.2
	github.com/hashicorp/consul/api v1.31.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/mbobakov/grpc-consul-resolver v1.5.3
	github.com/nicksnyder/go-i18n v1.10.3
	github.com/rabbitmq/amqp091-go v1.10.0
	go.opentelemetry.io/contrib/bridges/otelslog v0.11.0
	go.opentelemetry.io/otel/log v0.12.2 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0
	google.golang.org/grpc v1.72.2
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
