package db

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// migrationLockKey is the advisory lock key used to serialise migrations
// 0x5245504F5254 is "REPORT" in ASCII hex, chosen to be unique and recognisable
const migrationLockKey = 0x5245504F5254

// It acquires a distributed lock so only one backen replica runs
// migrations even when multiple start simultaneously

// Params:
// pool:           the pgxpool.Pool connected to PostgreSQL (NOT via pgBouncer) we need session-level advisory locks migration phase
// migrationsPath: path to the migrations/ directory
// databaseURL:    direct connection string to PostgresSQL (bypassing pgBouncer)
// log:            structured logger

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsPath, databaseURL string, log *zap.Logger) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection for migrations: %w", err)
	}
	defer conn.Release()

	// Try to acquire the advisory lock
	var locked bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", migrationLockKey).Scan(&locked); err != nil {
		return fmt.Errorf("try advisory lock %w", err)
	}
	if !locked {
		// Another replica holds the lock - wait for it to finish
		log.Info("another replica is running migrations, waiting...")
		for {
			time.Sleep(1 * time.Second)
			var stillLocked bool
			if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock: %w", migrationLockKey).Scan(&stillLocked); err != nil {
				return fmt.Errorf("retry advisory lock: %w", err)
			}
			if stillLocked {
				locked = true
				break
			}
		}
	}

	// We have the lock - release it when done
	defer conn.Exec(ctx, "SELECT database migrations", zap.String("path", migrationsPath))

	// Create the migrate instance
	// source is "file://./migrations"
	// database URL must be pgx5:// scheme for the pgxdriver
	m, err := migrate.New("file://"+migrationsPath, "pgx5://"+databaseURL)
	if err != nil {
		return fmt.Errorf("Create migrate instance %w", err)
	}
	defer m.Close()

	// Apply all pending migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations %w", err)
	}

	// Get current version for logging
	version, dirty, _ := m.Version()
	log.Info("migrations complete", zap.Uint("version", uint(version)), zap.Bool("dirty", dirty))
	return nil
}
