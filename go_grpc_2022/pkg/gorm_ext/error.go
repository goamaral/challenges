package gorm_ext

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type Error struct {
	Err        error
	WrappedErr error
}

func NewError(err error, wrappedErr error) Error {
	return Error{Err: err, WrappedErr: wrappedErr}
}
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.WrappedErr)
}
func (e Error) Unwrap() error {
	return e.WrappedErr
}

var (
	ErrUniqueViolation = errors.New("unique violation")
)

func ExtractError(err error) error {
	switch {
	case errors.Is(err, &pgconn.PgError{}):
		pgErr := err.(*pgconn.PgError)
		switch pgErr.SQLState() {
		case "23505": // unique violation
			return NewError(ErrUniqueViolation, pgErr)
		}
	}
	return err
}
