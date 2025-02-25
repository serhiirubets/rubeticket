package config

import (
	"github.com/joho/godotenv"
	"github.com/serhiirubets/rubeticket/pkg/convert"
	"log"
	"os"
)

type DbConfig struct {
	Dsn                             string
	MaxOpenConnections              int
	MaxIdleConnections              int
	MaxLifetimeConnectionsInMinutes int
}

type AuthConfig struct {
	Secret string
}

type Config struct {
	Db       DbConfig
	Auth     AuthConfig
	LogLevel string
	Env      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	maxOpenConnections := convert.StringToInt(os.Getenv("MaxOpenConnections"), 10)
	maxIdleConnections := convert.StringToInt(os.Getenv("MaxIdleConnections"), 10)
	maxLifetimeConnectionsInMinutes := convert.StringToInt(os.Getenv("MaxLifetimeConnectionsInMinutes"), 1)

	return &Config{
		Db: DbConfig{
			Dsn:                             os.Getenv("DSN"),
			MaxOpenConnections:              maxOpenConnections,
			MaxIdleConnections:              maxIdleConnections,
			MaxLifetimeConnectionsInMinutes: maxLifetimeConnectionsInMinutes,
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		LogLevel: os.Getenv("LOG_LEVEL"),
		Env:      os.Getenv("ENV"),
	}
}
