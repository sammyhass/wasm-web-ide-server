package projects

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/wasm"
)

type ProjectsService struct {
	repo        projectsRepo
	wasmService *wasm.WasmService
}

func NewProjectsService() *ProjectsService {
	return &ProjectsService{
		repo:        NewProjectsRepository(),
		wasmService: wasm.NewWasmService(),
	}
}

func (s *ProjectsService) CreateProject(
	name string,
	userID string,
) (model.ProjectView, error) {
	proj, err := s.repo.CreateProject(name, userID)
	if err != nil {
		return model.ProjectView{}, err
	}

	files, err := s.repo.CreateProjectFiles(proj)
	if err != nil {
		return model.ProjectView{}, err
	}

	return proj.ViewWithFiles(files), nil
}

func (s *ProjectsService) GetProjectsByUserID(userID string) ([]model.ProjectView, error) {
	return s.repo.GetProjectsByUserID(userID)
}

func (s *ProjectsService) GetProjectByID(userId, id string) (model.ProjectView, error) {
	return s.repo.GetProjectByID(userId, id)
}

func (s *ProjectsService) DeleteProjectByID(userId, id string) error {
	err := s.repo.DeleteProject(userId, id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteProjectFiles(userId, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectsService) CompileProjectWASM(
	userId string,
	projectId string,
) (string, error) {

	proj, err := s.repo.GetProjectByID(userId, projectId)

	if err != nil {
		return "", err
	}

	for _, file := range proj.Files {
		if file.Name == "main.go" {
			return s.wasmService.Compile(file.Content)
		}
	}

	return "", errors.New("main.go not found")
}

func (s *ProjectsService) UpdateProjectFiles(
	userId string,
	projectId string,
	files model.ProjectFiles,
) (
	model.ProjectFiles,
	error,
) {
	return s.repo.UpdateProjectFiles(userId, projectId, files)
}
