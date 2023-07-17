package shared

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

var (
	Pool *pgxpool.Pool

	DBConnection *sqlx.DB
)

func SetDBConnection(dbConn *sqlx.DB) {
	DBConnection = dbConn
}

func SetPGXPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("cannot assign nil pool")
	}

	Pool = newPool

	return nil
}
