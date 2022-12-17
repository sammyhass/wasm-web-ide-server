package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	*gorm.Model
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	Name      string
	UserID    string `gorm:"index"`
}

type ProjectView struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Name      string     `json:"name"`
	UserID    string     `json:"user_id"`
	Files     []FileView `json:"files"`
}

func (p *Project) View() ProjectView {
	return ProjectView{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
		UserID:    p.UserID,
	}
}

func (p *Project) ViewWithFiles(
	files ProjectFiles,
) ProjectView {
	return ProjectView{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
		UserID:    p.UserID,
		Files:     ProjectFilesToFileViews(files),
	}
}

func NewProject(
	name, userID string,
) Project {
	return Project{
		ID:     NewID(),
		Name:   name,
		UserID: userID,
	}
}
