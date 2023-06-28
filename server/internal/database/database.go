package database

import (
	"database/sql"
	"fmt"

	"erc-721-checks/internal/utils"

	_ "github.com/lib/pq"
)

const (
	maxIdleConnections    = 5
	maxOpenConnections    = 20
	connectionMaxLifetime = 100
)

var (
	DBInstance *sql.DB
	host       = utils.EnvHelper(utils.DBHost)
	port       = utils.EnvHelper(utils.DBPort)
	name       = utils.EnvHelper(utils.DBName)
	user       = utils.EnvHelper(utils.DBUser)
	password   = utils.EnvHelper(utils.DBPassword)
)

func InitDB() error {
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	var err error
	DBInstance, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	DBInstance.SetMaxIdleConns(maxIdleConnections)
	DBInstance.SetMaxOpenConns(maxOpenConnections)
	DBInstance.SetConnMaxLifetime(connectionMaxLifetime)

	if err := DBInstance.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}
