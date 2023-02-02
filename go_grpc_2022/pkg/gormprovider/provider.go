package gormprovider

import "github.com/jackc/pgconn"

const (
	UniqueViolationCode = "23505"
)

func IsUniqueViolationError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.SQLState() == UniqueViolationCode // unique_violation
}
