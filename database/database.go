package database

import (
	"database/sql"
)

func NewDatabase() *sql.DB {
	db, err := NewPostgres(Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "todos_db",
	})
	if err != nil {
		panic(err)
	}
	return db
}
