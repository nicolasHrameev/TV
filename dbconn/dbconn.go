package dbconn

import (
	"database/sql"
	"fmt"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

var db *sql.DB = nil

func GetDB() (*sql.DB, error) {
	if db == nil {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME)
		d, err := sql.Open("postgres", dbinfo)
		if err != nil {
			return nil, err
		}
		db = d
	}

	return db, nil
}
