package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Dbname   string
}

func NewPSQLStorage(cfg *DbConfig) (*sql.DB, error) {


	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error occured when connecting to db:", err)
	}

	return db, nil
}
