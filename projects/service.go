package projects

import (
	"strings"

	"github.com/sammyhass/web-ide/server/model"
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

	wasm, err := compileProject(mainFile)

	if err != nil {
		return "", err
	}
	go func() {
		wat, err := wasm2wat(wasm)

		if err != nil {
			return
		}

		r := strings.NewReader(wat)

		s.repo.uploadProjectWat(userId, projectId, r)
	}()

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
	return s.repo.uploadProjectSrcFiles(userId, projectId, files)
}

func (s *service) genProjectWatPresignedURL(userId, projectId string) (string, error) {
	return s.repo.genProjectWatPresignedURL(userId, projectId)
}

func (s *service) renameProject(userId, id, name string) (model.ProjectView, error) {
	return s.repo.renameProject(userId, id, name)
}
