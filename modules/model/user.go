package model

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID string `gorm:"primaryKey"`

	Email    string `gorm:"uniqueIndex"`
	Password string

	Projects []Project
}

type UserView struct {
	ID    string `json:"id"`
	Email string `json:"username"`
}

func (u *User) View() UserView {
	return UserView{
		ID:    u.ID,
		Email: u.Email,
	}
}
