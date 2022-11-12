package db

import (
	"fmt"
	"log"
	"time"

	"github.com/sammyhass/web-ide/server/modules/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
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

	sqlDb, err := conn.DB()
	if err != nil {
		log.Fatalf("[db] Error getting sql db: %v", err)
	}

	sqlDb.SetConnMaxLifetime(time.Minute * 4)

	db = conn
	log.Println("[db] Connected to database")
}

func GetConnection() *gorm.DB {
	if db == nil {
		Connect()
	}

	return db
}

func Close() {
	database, err := db.DB()
	if err != nil {
		log.Fatalf("[db] Error closing database: %v", err)
	}

	database.Close()
	log.Println("[db] Closed database connection")
}
