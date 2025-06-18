package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
)

func NewDBConn(cfg *config.Config) *pgxpool.Pool {
	connStr := cfg.Database.URL
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse database connection string: %v", err)
	}
	pgxConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	pgxConfig.MaxConnLifetime = cfg.Database.ConnMaxLifetime
	pgxConfig.MaxConnIdleTime = cfg.Database.ConnMaxIdleTime
	pgxConfig.HealthCheckPeriod = cfg.Database.HealthCheckPeriod
	pgxConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Printf("DB: Acquiring connection: %p", c)
		return true
	}
	pgxConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Printf("DB: Releasing connection: %p", c)
		return true
	}
	pgxConfig.BeforeClose = func(c *pgx.Conn) {
		log.Printf("DB: Closing connection: %p", c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected and pinged database!")
	return pool
}
