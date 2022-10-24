package projects

import "github.com/sammyhass/web-ide/server/modules/model"

type ProjectsService struct {
	repo *ProjectsRepository
}

func NewProjectsService() *ProjectsService {
	return &ProjectsService{
		repo: NewProjectsRepository(),
	}
}

func (s *ProjectsService) CreateProject(
	name string,
	description string,
	userID string,
) (model.Project, error) {
	return s.repo.CreateProject(name, description, userID)
}

func (s *ProjectsService) GetProjectsByUserID(userID string) ([]*model.Project, error) {
	return s.repo.GetProjectsByUserID(userID)
}

func (s *ProjectsService) GetProjectByID(id string) (*model.Project, error) {
	return s.repo.GetProjectByID(id)
}
