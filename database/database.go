package database

import (
	"database/sql"

	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPostgresConn(databaseURL string, databaseSchema string) *sql.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(databaseURL),
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": databaseSchema,
		})))

	return sqldb
}
