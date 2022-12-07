package model

import (
	"log"

	"github.com/sammyhass/web-ide/server/modules/db"
)

// Migrate is run on startup and should define the models (in this file)
func Migrate() {
	conn := db.GetConnection()

	if err := conn.AutoMigrate(&User{}, &Project{}); err != nil {
		log.Fatalf("Migration Failed: %v", err)
	}
}
