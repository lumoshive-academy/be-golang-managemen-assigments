package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	conntSTR := "user=postgres password=postgres dbname=assignment sslmode=disable host=192.168.1.25"
	db, err := sql.Open("postgres", conntSTR)
	return db, err
}
