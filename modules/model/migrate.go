package model

import (
	"github.com/sammyhass/web-ide/server/modules/db"
)

// Migrate is run on startup and should define the models (in this file)
func Migrate() {
	conn := db.GetConnection()

	conn.AutoMigrate(&User{}, &Project{})
}
