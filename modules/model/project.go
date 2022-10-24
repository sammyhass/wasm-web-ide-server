package model

import "gorm.io/gorm"

type Project struct {
	*gorm.Model
	ID          string `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	UserID string `json:"user_id"`
}

type ProjectView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *Project) View() ProjectView {
	return ProjectView{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}
