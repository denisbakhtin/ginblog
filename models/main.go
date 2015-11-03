package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func SetDB(connection string) {
	db = sqlx.MustConnect("postgres", connection)
}

func GetDB() *sqlx.DB {
	return db
}
