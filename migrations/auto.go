package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, connectErr := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if connectErr != nil {
		panic(connectErr)
	}

	migrateErr := db.Migrator().AutoMigrate(&users.User{}, &file.File{})
	if migrateErr != nil {
		fmt.Println(migrateErr.Error())
		return
	}
}
