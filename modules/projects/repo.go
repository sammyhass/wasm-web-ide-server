package projects

import (
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/s3"
	"gorm.io/gorm"
)

type projectsRepo interface {
	createProject(
		name string,
		userID string,
	) (model.Project, error)

	createProjectFiles(project model.Project) (model.ProjectFiles, error)

	getProjectsByUserID(userID string) ([]model.ProjectView, error)
	getProjectByID(userId string, id string) (model.ProjectView, error)

	deleteProject(userId string, id string) error
	deleteProjectFiles(userId string, id string) error

	updateProjectSrcFiles(userId string, id string, files model.ProjectFiles) (model.ProjectFiles, error)

	uploadProjectWasm(userId string, id string, f io.Reader) error

	getProjectWasmPresignedURL(userId string, id string) (string, error)
}

type Repository struct {
	db *gorm.DB
	s3 *s3.Service
}

func NewProjectsRepository() *Repository {
	return &Repository{
		db: db.GetConnection(),
		s3: s3.NewS3Service(),
	}
}

/*
createProject creates a new project in the database
*/
func (r *Repository) createProject(
	name string,
	userID string,
) (model.Project, error) {

	proj := model.NewProject(
		name,
		userID,
	)

	err := r.db.Create(&proj).Error

	if err != nil {
		return model.Project{}, err
	}

	return proj, nil
}

func (r *Repository) createProjectFiles(project model.Project) (model.ProjectFiles, error) {
	goMod := model.DefaultGoMod(project.Name)
	files := model.DefaultFiles

	files["go.mod"] = goMod

	srcDir := fmt.Sprintf("%s/%s/src", project.UserID, project.ID)

	if err := r.s3.UploadFiles(srcDir, files); err != nil {
		return nil, err
	}

	return model.DefaultFiles, nil
}

/*
getProjectsByUserID returns all projects for a given user (without files)
*/
func (r *Repository) getProjectsByUserID(userID string) ([]model.ProjectView, error) {
	var projects []*model.Project

	err := r.db.Where("user_id = ?", userID).Find(&projects).Error
	if err != nil {
		return nil, err
	}

	out := make([]model.ProjectView, len(projects))
	for i := range projects {
		out[i] = projects[i].View()
	}

	return out, nil
}

/*
getProjectByID returns a project for a given user with the given id returning the database record and the files in s3
*/
func (r *Repository) getProjectByID(userId string, id string) (model.ProjectView, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return model.ProjectView{}, err
	}

	if project.UserID != userId {
		return model.ProjectView{}, errors.New("project not found")
	}

	files, err := r.s3.GetFiles(fmt.Sprintf("%s/%s/src", userId, id))

	if err != nil {
		return model.ProjectView{}, err
	}

	return project.ViewWithFiles(files), nil
}

/*
DeleteProject deletes a project from the database
*/
func (r *Repository) deleteProject(userId string, id string) error {

	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return err
	}

	if project.UserID != userId {
		return errors.New("project not found")
	}

	dbErr := r.db.Delete(&project).Error
	if dbErr != nil {
		return dbErr
	}

	return nil
}

/*
DeleteDir deletes a project stored in s3
*/
func (r *Repository) deleteProjectFiles(userId string, id string) error {
	dir := fmt.Sprintf("%s/%s", userId, id)
	s3Err := r.s3.DeleteDir(dir)
	if s3Err != nil {
		return s3Err
	}

	return nil
}

/*
updateProjectFiles updates the files for a given project in s3
*/
func (r *Repository) updateProjectSrcFiles(userId string, id string, files model.ProjectFiles) (
	model.ProjectFiles,
	error,
) {
	srcDir := fmt.Sprintf("%s/%s/src", userId, id)
	err := r.s3.UploadFiles(srcDir, files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *Repository) uploadProjectWasm(userId string, id string, file io.Reader) error {
	wasmDir := fmt.Sprintf("%s/%s/build", userId, id)
	_, err := r.s3.Upload(wasmDir, "main.wasm", file)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) getProjectWasmPresignedURL(userId string, id string) (string, error) {
	wasmDir := fmt.Sprintf("%s/%s/build", userId, id)
	url, err := r.s3.GenPresignedURL(path.Join(wasmDir, "main.wasm"), time.Hour*24*7)
	if err != nil {
		return "", err
	}

	return url, nil
}
