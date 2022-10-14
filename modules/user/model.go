package user

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`

	Password string
	Salt     string

	Projects []Project
}

type Project struct {
	*gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`

	UserID string
}
