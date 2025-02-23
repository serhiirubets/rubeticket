package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
