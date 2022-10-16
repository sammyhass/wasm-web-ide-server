package migrations

import (
	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/user"
)

func Migrate() {
	conn := db.GetConnection()

	conn.AutoMigrate(&user.User{})
}
