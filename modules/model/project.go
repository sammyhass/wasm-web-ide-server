package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	*gorm.Model
	ID          string `gorm:"primaryKey" json:"id"`
	Name        string
	Description string

	UserID string

	WasmBinPath string
}

type ProjectView struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	UserID      string     `json:"user_id"`
	WasmBinPath string     `json:"wasm_path"`
	Files       []FileView `json:"files"`
}

func (p *Project) View() ProjectView {
	return ProjectView{
		ID:          p.ID,
		CreatedAt:   p.CreatedAt,
		Name:        p.Name,
		UserID:      p.UserID,
		Description: p.Description,
		WasmBinPath: p.WasmBinPath,
	}
}

func (p *Project) ViewWithFiles(files []FileView) ProjectView {
	return ProjectView{
		ID:          p.ID,
		CreatedAt:   p.CreatedAt,
		Name:        p.Name,
		UserID:      p.UserID,
		Description: p.Description,
		WasmBinPath: p.WasmBinPath,
		Files:       files,
	}
}

func NewProject(
	name, description, userID string,
) Project {
	return Project{
		ID:          NewID(),
		Name:        name,
		Description: description,
		UserID:      userID,
	}
}
