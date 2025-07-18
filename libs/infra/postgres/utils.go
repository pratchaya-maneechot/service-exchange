package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToUUID(str string) pgtype.UUID {
	if str == "" {
		return pgtype.UUID{Valid: false}
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}
}

func FromUUID(val pgtype.UUID) string {
	if !val.Valid {
		return ""
	}
	return uuid.UUID(val.Bytes).String()
}

func ToTime(val pgtype.Timestamptz) *time.Time {
	if !val.Valid {
		return nil
	}
	return &val.Time
}

func ToTimestamp(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

func ToText(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: str,
		Valid:  true,
	}
}

func FromText(val pgtype.Text) string {
	if !val.Valid {
		return ""
	}
	return val.String
}

func ToInt4(val int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(val),
		Valid: true,
	}
}

func FromInt4(val pgtype.Int4) int {
	if !val.Valid {
		return 0
	}
	return int(val.Int32)
}

func ToBool(val bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  val,
		Valid: true,
	}
}

func FromBool(val pgtype.Bool) bool {
	if !val.Valid {
		return false
	}
	return val.Bool
}
