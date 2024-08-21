package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type Tracer struct{}

// TraceQueryStart is called at the beginning of a Query, QueryRow, or Exec call.
func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	startTime := time.Now()
	slog.Debug("cases.store.query_started",
		slog.String("sql", data.SQL),
		slog.Any("args", data.Args),
		slog.Time("start_time", startTime),
	)

	// Store start time in context for later use in TraceQueryEnd
	ctx = context.WithValue(ctx, "start_time", startTime)
	return ctx
}

// TraceQueryEnd is called at the end of a Query, QueryRow, or Exec call.
func (t *Tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	startTime, ok := ctx.Value("start_time").(time.Time)
	if !ok {
		return // If start_time is not found in the context, skip logging
	}
	duration := time.Since(startTime)

	if data.Err != nil {
		slog.Error("cases.store.query_error",
			slog.Duration("duration", duration),
			slog.String("error", data.Err.Error()),
		)
	} else {
		slog.Debug("cases.store.query_completed",
			slog.Duration("duration", duration),
			slog.Int64("rows_affected", data.CommandTag.RowsAffected()),
		)
	}
}
