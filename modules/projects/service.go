package projects

import (
	"github.com/sammyhass/web-ide/server/modules/model"
)

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
) (model.ProjectView, error) {
	// Create project in db
	pv, err := s.repo.CreateProject(name, description, userID)
	if err != nil {
		return model.ProjectView{}, err
	}

	return pv, nil
}

func (s *ProjectsService) GetProjectsByUserID(userID string) ([]model.ProjectView, error) {
	return s.repo.GetProjectsByUserID(userID)
}

func (s *ProjectsService) GetProjectByID(userId, id string) (model.ProjectView, error) {
	return s.repo.GetProjectByID(userId, id)
}
