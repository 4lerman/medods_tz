package app

import (
	"context"
	"database/sql"

	"github.com/4lerman/medods_tz/api"
	configs "github.com/4lerman/medods_tz/internal/config"
	"github.com/4lerman/medods_tz/pkg/db"
	"github.com/4lerman/medods_tz/pkg/log"
	"go.uber.org/zap"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())


	db, err := db.NewPSQLStorage(&db.DbConfig{
		Host:     configs.Envs.DBAddress,
		User:     configs.Envs.DBUser,
		Port:     configs.Envs.DBPort,
		Dbname:   configs.Envs.DBName,
		Password: configs.Envs.DBPassword,
	})

	if err != nil {
		logger.Fatal("Db init error:", zap.Error(err))
	}

	defer db.Close()

	initStorage(db)

	server := api.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		logger.Fatal("Error when running server", zap.Error(err))
	}
}

func initStorage(db *sql.DB) {
	logger := log.LoggerFromContext(context.Background())

	err := db.Ping()
	if err != nil {
		logger.Fatal("DB connection error:", zap.Error(err))
	}

	logger.Info("DB connected successfully")
}
