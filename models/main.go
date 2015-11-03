package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgresql driver, don't remove
)

var db *sqlx.DB

//SetDB establishes connection to database and saves its handler into db *sqlx.DB
func SetDB(connection string) {
	db = sqlx.MustConnect("postgres", connection)
}

//GetDB returns database handler
func GetDB() *sqlx.DB {
	return db
}
