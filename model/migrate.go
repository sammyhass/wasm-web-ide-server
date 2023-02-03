package model

import (
	"log"

	"github.com/sammyhass/web-ide/server/db"
)

// Migrate performs a database migration
func Migrate() {

	conn := db.GetConnection()

	if err := conn.AutoMigrate(&User{}, &Project{}); err != nil {
		log.Fatalf("Migration Failed: %v", err)
	}
}
