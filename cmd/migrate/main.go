package main

import (
	"log"
	"os"

	configs "github.com/4lerman/medods_tz/internal/config"
	"github.com/4lerman/medods_tz/pkg/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := db.NewPSQLStorage(&db.DbConfig{
		Host:     configs.Envs.DBAddress,
		User:     configs.Envs.DBUser,
		Port:     configs.Envs.DBPort,
		Dbname:   configs.Envs.DBName,
		Password: configs.Envs.DBPassword,
	})

	if err != nil {
		log.Fatal("Db init error:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal("Migration error: ", err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
