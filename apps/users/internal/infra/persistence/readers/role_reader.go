package readers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
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

	var dRoles []role.Role
	for _, v := range roles {
		id := uint(v.ID)
		dRoles = append(dRoles, role.Role{
			ID:          &id,
			Name:        role.RoleName(v.Name),
			Description: *v.Description,
		})
	}
	return dRoles, nil
}
