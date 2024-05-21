package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "usagestats.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
