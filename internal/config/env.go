package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost         string
	Port               string
	DBUser             string
	DBPassword         string
	DBAddress          string
	DBPort             int64
	DBName             string
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExp     int64
	RefreshTokenExp    int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:         getEnv("PUBLIC_HOST", "http://localhost"),
		Port:               getEnv("PORT", "8080"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:          getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnvAsInt("DB_PORT", 5432),
		DBName:             getEnv("DB_NAME", "ecom"),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "not-so-secret-now-is-it?"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "not-so-secret-now-is-it?"),
		AccessTokenExp:     getEnvAsInt("ACCESS_TOKEN_EXP", 60*5),
		RefreshTokenExp:    getEnvAsInt("REFRESH_TOKEN_EXP", 60*60*5),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
