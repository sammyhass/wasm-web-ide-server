package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sammyhass/web-ide/server/modules/env"
)

var db *gorm.DB

func createConnection() *gorm.DB {
	if db != nil {
		return db
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.Env.POSTGRES_HOST,
		env.Env.POSTGRES_USER,
		env.Env.POSTGRES_PASSWORD,
		env.Env.POSTGRES_DB,
		env.Env.POSTGRES_PORT,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db = conn

	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		db = createConnection()
	}
	return db
}
