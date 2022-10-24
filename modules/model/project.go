package model

import "gorm.io/gorm"

type Project struct {
	*gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`

	UserID string
}
