package projects

import (
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/wasm"
)

type Service struct {
	repo        projectsRepo
	wasmService *wasm.Service
}

func NewService() *Service {
	return &Service{
		repo:        NewProjectsRepository(),
		wasmService: wasm.NewWasmService(),
	}
}

func (s *Service) CreateProject(
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

func (s *Service) GetProjectsByUserID(userID string) ([]model.ProjectView, error) {
	return s.repo.GetProjectsByUserID(userID)
}

func (s *Service) GetProjectByID(userId, id string) (model.ProjectView, error) {
	return s.repo.GetProjectByID(userId, id)
}

func (s *Service) DeleteProjectByID(userId, id string) error {
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

func (s *Service) CompileProjectWASM(
	userId string,
	projectId string,
) (string, error) {

	proj, err := s.repo.GetProjectByID(userId, projectId)

	if err != nil {
		return "", err
	}

	mainFile, err := model.GetFileContent(proj.Files, "main.go")
	if err != nil {
		return "", err
	}

	modFile, err := model.GetFileContent(proj.Files, "go.mod")
	if err != nil {
		return "", err
	}

	return s.wasmService.Compile(mainFile, modFile)
}

func (s *Service) UpdateProjectFiles(
	userId string,
	projectId string,
	files model.ProjectFiles,
) (
	model.ProjectFiles,
	error,
) {
	return s.repo.UpdateProjectFiles(userId, projectId, files)
}
