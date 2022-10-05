package projects

import "gorm.io/gorm"

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (pr *ProjectRepository) Create(name string) (*Project, error) {
	proj := NewProject(name)

	tx := pr.db.Create(proj)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return proj, nil
}
