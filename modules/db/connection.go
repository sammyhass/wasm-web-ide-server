package db

import (
	"fmt"
	"log"

	"github.com/sammyhass/web-ide/server/modules/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {

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
		log.Fatalf("[db] Error connecting to database: %v", err)
	}

	db = conn
	log.Println("[db] Connected to database")
}

func GetConnection() *gorm.DB {
	return db
}
