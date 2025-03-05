package tests

import (
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
)

type TestEnv struct {
	DB     *db.Db
	Logger log.ILogger
	Conf   *config.Config
}

func SetupTestEnv() *TestEnv {
	conf := config.LoadConfig()
	dbInstance := db.NewDb(conf)
	logger := log.NewLogrusLogger(conf.LogLevel)

	return &TestEnv{
		DB:     dbInstance,
		Logger: logger,
		Conf:   conf,
	}
}
