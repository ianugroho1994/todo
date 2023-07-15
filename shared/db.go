package shared

import "github.com/jmoiron/sqlx"

var DBConnection *sqlx.DB

func SetDBConnection(dbConn *sqlx.DB) {
	DBConnection = dbConn
}
