package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID string `gorm:"primaryKey"`

	Username string `gorm:"uniqueIndex"`
	Password string

	Projects []Project
}

func NewID() string {
	return uuid.New().String()
}
