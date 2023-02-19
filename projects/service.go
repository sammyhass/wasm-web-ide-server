package projects

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/sammyhass/web-ide/server/model"
	"github.com/sammyhass/web-ide/server/wasm"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: newRepository(),
	}
}

func (s *Service) CreateProject(
	name string,
	userId string,
	language model.ProjectLanguage,
) (model.ProjectView, error) {
	proj, err := s.repo.createProject(name, userId, language)
	if err != nil {
		return model.ProjectView{}, err
	}

	files, err := s.repo.createProjectFiles(userId, userId, language)
	if err != nil {
		return model.ProjectView{}, err
	}

	return proj.ViewWithFiles(files), nil
}

func (s *Service) GetProjectsByUserID(userId string) ([]model.ProjectView, error) {
	return s.repo.getProjectsByUserID(userId)
}

func (s *Service) GetProjectByID(userId, id string) (model.ProjectView, error) {
	return s.repo.getProjectByID(userId, id)
}

func (s *Service) DeleteProjectByID(userId, id string) error {
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

func (s *Service) CompileProjectWASM(
	userId string,
	projectId string,
) (string, error) {

	proj, err := s.repo.getProjectByID(userId, projectId)
	if err != nil {
		return "", err
	}

	var fname string
	switch proj.Language {
	case model.LanguageAssemblyScript.String():
		fname = "main.ts"
	case model.LanguageGo.String():
		fname = "main.go"
	default:
		return "", errors.New("unsupported language")
	}

	mainFile, err := model.GetFileContent(proj.Files, fname)
	if err != nil {
		return "", err
	}

	res, err := wasm.Compile(
		model.GetProjectLanguage(proj.Language),
		mainFile,
		wasm.CompileOpts{
			GenWat: true,
		},
	)
	if err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	var uploadErrs []error

	wg.Add(2)
	wasmReader := bytes.NewReader(res.Wasm)

	go func() {
		defer wg.Done()

		if err = s.repo.uploadProjectWasm(userId, projectId, wasmReader); err != nil {
			uploadErrs = append(uploadErrs, err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err != nil {
			uploadErrs = append(uploadErrs, err)
			return
		}
		watReader := strings.NewReader(res.Wat)
		if err = s.repo.uploadProjectWat(userId, projectId, watReader); err != nil {
			uploadErrs = append(uploadErrs, err)
			return
		}
	}()

	wg.Wait()

	if len(uploadErrs) > 0 {
		return "", uploadErrs[0]
	}

	return s.repo.genProjectWasmPresignedURL(userId, projectId)
}

func (s *Service) UpdateProjectFiles(
	userId string,
	projectId string,
	files model.ProjectFiles,
) (
	model.ProjectFiles,
	error,
) {
	return s.repo.uploadProjectSrcFiles(userId, projectId, files)
}

func (s *Service) GenProjectWatPresignedURL(userId, projectId string) (string, error) {
	return s.repo.genProjectWatPresignedURL(userId, projectId)
}

func (s *Service) RenameProject(userId, id, name string) (model.ProjectView, error) {
	return s.repo.renameProject(userId, id, name)
}

func (s *Service) ToggleForkable(userId, id string, share bool) (sharecode string, err error) {
	return s.repo.toggleSharing(userId, id, share)
}

func (s *Service) ForkProject(userId, sharecode string) (model.ProjectView, error) {
	sharedProject, err := s.repo.getProjectByShareCode(sharecode)
	if err != nil {
		return model.ProjectView{}, err
	}

	newProj, err := s.repo.createProject(
		fmt.Sprintf("%s (fork)", sharedProject.Name),
		userId,
		model.GetProjectLanguage(sharedProject.Language),
	)

	if err != nil {
		return model.ProjectView{}, err
	}

	files, err := s.repo.createProjectFilesWith(
		userId,
		newProj.ID,
		model.FileViewsToProjectFiles(sharedProject.Files),
	)
	if err != nil {
		return model.ProjectView{}, err
	}

	return newProj.ViewWithFiles(files), nil
}

func (s *Service) GetSharedProject(sharecode string) (model.ProjectView, error) {
	return s.repo.getProjectByShareCode(sharecode)
}
