package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
)

type DBPool struct {
	Pool *pgxpool.Pool
}

func NewDBConn(cfg *config.Config, log *slog.Logger) (*DBPool, error) {
	connStr := cfg.Database.URL
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	pgxConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	pgxConfig.MaxConnLifetime = cfg.Database.ConnMaxLifetime
	pgxConfig.MaxConnIdleTime = cfg.Database.ConnMaxIdleTime
	pgxConfig.HealthCheckPeriod = cfg.Database.HealthCheckPeriod

	pgxConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Debug("DB: Acquiring connection", "conn_ptr", fmt.Sprintf("%p", c))
		return true
	}
	pgxConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Debug("DB: Releasing connection", "conn_ptr", fmt.Sprintf("%p", c))
		return true
	}
	pgxConfig.BeforeClose = func(c *pgx.Conn) {
		log.Debug("DB: Closing connection", "conn_ptr", fmt.Sprintf("%p", c))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected and pinged database!")
	return &DBPool{
		Pool: pool,
	}, nil
}

func (d *DBPool) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
