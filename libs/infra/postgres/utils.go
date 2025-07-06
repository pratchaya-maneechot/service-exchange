package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func EncodeUID(str string) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(uuid.MustParse(str)),
		Valid: true,
	}
}
