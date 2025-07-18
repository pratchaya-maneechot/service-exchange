package postgres

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	lp "github.com/pratchaya-maneechot/service-exchange/libs/infra/postgres"
)

func NewDBConn(cfg *config.Config, log *slog.Logger) (*lp.DBPool, error) {
	connStr := cfg.Database.URL
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}
	conn, err := lp.NewConnection(pgxConfig, log)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
