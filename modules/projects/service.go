package projects

import (
	"github.com/sammyhass/web-ide/server/modules/model"
)

type service struct {
	repo projectsRepo
}

func newService() *service {
	return &service{
		repo: NewProjectsRepository(),
	}
}

func (s *service) createProject(
	name string,
	userID string,
) (model.ProjectView, error) {
	proj, err := s.repo.createProject(name, userID)
	if err != nil {
		return model.ProjectView{}, err
	}

	files, err := s.repo.createProjectFiles(proj)
	if err != nil {
		return model.ProjectView{}, err
	}

	return proj.ViewWithFiles(files), nil
}

func (s *service) getProjectsByUserID(userID string) ([]model.ProjectView, error) {
	return s.repo.getProjectsByUserID(userID)
}

func (s *service) getProjectByID(userId, id string) (model.ProjectView, error) {
	return s.repo.getProjectByID(userId, id)
}

func (s *service) deleteProjectByID(userId, id string) error {
	err := s.repo.deleteProject(userId, id)
	if err != nil {
		return err
	}

	err = s.repo.deleteProjectFiles(userId, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) compileProjectWASM(
	userId string,
	projectId string,
) (string, error) {

	proj, err := s.repo.getProjectByID(userId, projectId)

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

	wasm, err := compileProject(mainFile, modFile)

	if err != nil {
		return "", err
	}

	if err = s.repo.uploadProjectWasm(userId, projectId, wasm); err != nil {
		return "", err
	}

	return s.repo.genProjectWasmPresignedURL(userId, projectId)
}

func (s *service) updateProjectFiles(
	userId string,
	projectId string,
	files model.ProjectFiles,
) (
	model.ProjectFiles,
	error,
) {
	return s.repo.updateProjectSrcFiles(userId, projectId, files)
}
