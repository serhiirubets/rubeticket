package config

import (
	"github.com/joho/godotenv"
	"github.com/serhiirubets/rubeticket/internal/pkg/convert"
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

type AppConfig struct {
	Port string
	Host string
}

type Config struct {
	Db       DbConfig
	Auth     AuthConfig
	LogLevel string
	Env      string
	App      AppConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	maxOpenConnections := convert.StringToInt(os.Getenv("MAX_OPEN_CONNECTIONS"), 10)
	maxIdleConnections := convert.StringToInt(os.Getenv("MAX_IDLE_CONNECTIONS"), 10)
	maxLifetimeConnectionsInMinutes := convert.StringToInt(os.Getenv("MAX_LIFE_TIME_CONNECTIONS_IN_MINUTES"), 1)

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
		App: AppConfig{
			Port: os.Getenv("PORT"),
			Host: os.Getenv("HOST"),
		},
	}
}
