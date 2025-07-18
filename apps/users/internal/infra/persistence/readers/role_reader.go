package readers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
	lp "github.com/pratchaya-maneechot/service-exchange/libs/infra/postgres"
	"github.com/pratchaya-maneechot/service-exchange/libs/utils"
)

type roleReader struct {
	db   *db.Queries
	pool *pgxpool.Pool
}

func NewPostgresRoleReader(dbPool *lp.DBPool) role.RoleReader {
	return &roleReader{
		db:   db.New(dbPool.Pool),
		pool: dbPool.Pool,
	}
}

func (r *roleReader) GetAllRoles(ctx context.Context) ([]role.Role, error) {
	roles, err := r.db.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	return utils.ArrayMap(roles, func(dr db.Role) role.Role { return *role.NewRoleFromRepository(uint(dr.ID), dr.Name, dr.Description) }), nil
}
