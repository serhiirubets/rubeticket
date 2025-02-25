package db

import (
	configs "github.com/serhiirubets/rubeticket/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(conf.Db.MaxOpenConnections)

	sqlDB.SetMaxIdleConns(conf.Db.MaxIdleConnections)

	sqlDB.SetConnMaxLifetime(time.Duration(conf.Db.MaxLifetimeConnectionsInMinutes) * time.Minute)

	return &Db{db}
}
