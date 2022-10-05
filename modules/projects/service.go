package projects

import "gorm.io/gorm"

type ProjectsService struct {
	repo *ProjectRepository
}

func NewService(db *gorm.DB) ProjectsService {
	return ProjectsService{
		repo: NewProjectRepository(db),
	}
}

func (pr *ProjectsService) GetProjects() []Project {
	return nil
}

func (pr *ProjectsService) CreateProject(name string) (*Project, error) {
	return pr.repo.Create(name)
}
