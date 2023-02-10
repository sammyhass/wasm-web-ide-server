package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectLanguage is the language chosen to be compiled to WASM for a given project
type ProjectLanguage int

var langs = [...]string{
	"Go",
	"AssemblyScript",
}

const (
	LanguageGo ProjectLanguage = iota
	LanguageAssemblyScript
)

func GetProjectLanguage(name string) ProjectLanguage {
	for i, lang := range langs {
		if lang == name {
			return ProjectLanguage(i)
		}
	}
	return LanguageGo
}

func (l ProjectLanguage) String() string {
	return langs[l]
}

type Project struct {
	*gorm.Model
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	Name      string
	UserID    string          `gorm:"index"`
	Language  ProjectLanguage `gorm:"default:0"`
}

type ProjectView struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Name      string     `json:"name"`
	UserID    string     `json:"user_id"`
	Files     []FileView `json:"files"`
	WasmPath  string     `json:"wasm_path"`
	Language  string     `json:"language"`
}

func (p *Project) View() ProjectView {
	return ProjectView{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
		UserID:    p.UserID,
		Language:  p.Language.String(),
	}
}

func (p *Project) ViewWith(
	wasmPath string,
	files ProjectFiles,
) ProjectView {
	return ProjectView{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		Name:      p.Name,
		UserID:    p.UserID,
		WasmPath:  wasmPath,
		Files:     ProjectFilesToFileViews(files),
		Language:  p.Language.String(),
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
		Language:  p.Language.String(),
		Files:     ProjectFilesToFileViews(files),
	}
}

func NewProject(
	name, userID string,
	language ProjectLanguage,
) Project {
	return Project{
		ID:       NewID(),
		Name:     name,
		UserID:   userID,
		Language: language,
	}
}
