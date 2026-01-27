package store

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
