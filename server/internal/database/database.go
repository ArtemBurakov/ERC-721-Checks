package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host                  = "localhost"
	port                  = 5432
	user                  = "postgres"
	password              = "1111"
	databaseName          = "checks"
	maxIdleConnections    = 5
	maxOpenConnections    = 20
	connectionMaxLifetime = 100
)

var db *sql.DB

func InitDB() error {
	var err error
	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, databaseName)
	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connectionMaxLifetime)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
