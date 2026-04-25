package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	// is the connection string like: postgres://user:password@host:port/database
	// in production via pgbouncer using :6432
	DNS string
	// baypass pgbouncer directly connect to postgres on :5432
	DirectDNS string
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DNS)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 3

	// close conns to prevent issues with stale connections that postgres may have
	// cleaned up
	poolConfig.MaxConnLifetime = 30 * time.Minute

	// close the connection unused longer than 5 min
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	poolConfig.HealthCheckPeriod = 30 * time.Minute

	// CRITICAL: pgBouncer transaction mode
	// pgbouncer does not support prepared statement (extended query protocol)
	// without this, you get cryptic errors like "unexpected message type"
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
