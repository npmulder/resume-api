// Package database provides database functionality for the application.
package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracedPool is a wrapper around pgxpool.Pool that adds tracing to database operations.
type TracedPool struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
}

// NewTracedPool creates a new TracedPool.
func NewTracedPool(pool *pgxpool.Pool) *TracedPool {
	return &TracedPool{
		pool:   pool,
		tracer: otel.GetTracerProvider().Tracer("database"),
	}
}

// Acquire acquires a connection from the pool with tracing.
func (tp *TracedPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	ctx, span := tp.tracer.Start(ctx, "db.acquire")
	defer span.End()

	conn, err := tp.pool.Acquire(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	return conn, nil
}

// Begin begins a transaction with tracing.
func (tp *TracedPool) Begin(ctx context.Context) (pgx.Tx, error) {
	ctx, span := tp.tracer.Start(ctx, "db.begin_transaction")
	defer span.End()

	tx, err := tp.pool.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	return tx, nil
}

// Exec executes a query with tracing.
func (tp *TracedPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	ctx, span := tp.tracer.Start(ctx, "db.exec")
	defer span.End()

	span.SetAttributes(attribute.String("db.statement", sql))

	result, err := tp.pool.Exec(ctx, sql, arguments...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return result, err
	}

	span.SetAttributes(attribute.Int64("db.rows_affected", int64(result.RowsAffected())))
	return result, nil
}

// Query executes a query with tracing.
func (tp *TracedPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	ctx, span := tp.tracer.Start(ctx, "db.query")
	defer span.End()

	span.SetAttributes(attribute.String("db.statement", sql))

	rows, err := tp.pool.Query(ctx, sql, args...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	return rows, nil
}

// QueryRow executes a query that returns a single row with tracing.
func (tp *TracedPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	ctx, span := tp.tracer.Start(ctx, "db.query_row")
	defer span.End()

	span.SetAttributes(attribute.String("db.statement", sql))

	return tp.pool.QueryRow(ctx, sql, args...)
}

// Close closes the pool.
func (tp *TracedPool) Close() {
	tp.pool.Close()
}

// Ping pings the database.
func (tp *TracedPool) Ping(ctx context.Context) error {
	ctx, span := tp.tracer.Start(ctx, "db.ping")
	defer span.End()

	err := tp.pool.Ping(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	return err
}

// Stat returns statistics about the connection pool.
func (tp *TracedPool) Stat() *pgxpool.Stat {
	return tp.pool.Stat()
}

// Pool returns the underlying pgxpool.Pool.
func (tp *TracedPool) Pool() *pgxpool.Pool {
	return tp.pool
}
