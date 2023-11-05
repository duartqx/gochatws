package repositories

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func GetDbConnection(dbfile string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	// Enforces FK constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	return db, nil
}
