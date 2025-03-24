package db

import (
	"log"
	"os"
	"time"

	configs "github.com/serhiirubets/rubeticket/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) IDb {
	var gormLogger logger.Interface

	if conf.Env == "dev" {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // For slow request
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		// For other env log only errors
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		panic(err)
	}

	pgDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	pgDb.SetMaxOpenConns(conf.Db.MaxOpenConnections)

	pgDb.SetMaxIdleConns(conf.Db.MaxIdleConnections)

	pgDb.SetConnMaxLifetime(time.Duration(conf.Db.MaxLifetimeConnectionsInMinutes) * time.Minute)

	return &Db{db}
}
