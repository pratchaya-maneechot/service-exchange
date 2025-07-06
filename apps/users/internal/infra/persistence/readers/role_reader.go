package readers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
	"github.com/pratchaya-maneechot/service-exchange/libs/utils"
)

// PostgresRoleReader implements the role.RoleReader interface.
type roleReader struct {
	db   *db.Queries
	pool *pgxpool.Pool // Keep pool for transaction management
}

// NewPostgresRoleReader creates a new PostgresRoleReader.
func NewPostgresRoleReader(dbPool *postgres.DBPool) role.RoleReader {
	return &roleReader{
		db:   db.New(dbPool.Pool),
		pool: dbPool.Pool,
	}
}

// GetAllRoles implements role.RoleReader.
func (r *roleReader) GetAllRoles(ctx context.Context) ([]role.Role, error) {
	roles, err := r.db.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	return utils.ArrayMap(roles, func(dr db.Role) role.Role { return *role.NewRoleFromRepository(uint(dr.ID), dr.Name, dr.Description) }), nil
}
