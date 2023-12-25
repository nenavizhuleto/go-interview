package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DBNAME = "app.db"

var DB *sqlx.DB

func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", DBNAME)
	return db, err
}
