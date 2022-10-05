package migrations

import (
	"github.com/sammyhass/web-ide/server/modules/projects"
	"gorm.io/gorm"
)

var models []interface{} = []interface{}{
	projects.Project{},
}

func Migrate(
	db *gorm.DB,
) {
	db.AutoMigrate(models...)
}
