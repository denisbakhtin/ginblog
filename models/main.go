package models

import (
	"fmt"

	"github.com/denisbakhtin/ginbasic/system"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func SetDB(config *system.Config) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)

	db = sqlx.MustConnect("postgres", connectionString)
}

func GetDB() *sqlx.DB {
	return db
}
