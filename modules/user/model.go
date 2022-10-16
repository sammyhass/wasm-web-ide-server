package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`

	Password string

	// Projects []Project
}

func NewID() string {
	return uuid.New().String()
}

// type Project struct {
// 	*gorm.Model
// 	ID   string `gorm:"primaryKey"`
// 	Name string `gorm:"uniqueIndex"`

// 	ExecutableSrc string `gorm:"uniqueIndex"` // Path to wasm file

// 	Files []File

// 	UserID string
// }

// type File struct {
// 	*gorm.Model
// 	ID   string `gorm:"primaryKey"`
// 	Name string `gorm:"uniqueIndex"`

// 	Content string

// 	ProjectID string
// }
