package projects

import (
	"errors"
	"io"
	"path"
	"time"

	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/s3"
	"gorm.io/gorm"
)

// getProjectDir returns the path to the project directory as stored in the s3 bucket
func getProjectDir(userId string, id string) string {
	return path.Join(userId, id)
}

// getProjectSrcDir returns the path to the project source directory as stored in the s3 bucket
func getProjectSrcDir(userId string, id string) string {
	return path.Join(getProjectDir(userId, id), "src")
}

// getProjectWasmDir returns the path to the project wasm directory as stored in the s3 bucket
func getProjectWasmDir(userId string, id string) string {
	return path.Join(getProjectDir(userId, id), "build")
}

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

	uploadProjectWasm(userId string, id string, r io.Reader) error

	genSrcPresignedURL(userId string, id string, filename string) (string, error)

	genProjectWasmPresignedURL(userId string, id string) (string, error)
}

type repository struct {
	db *gorm.DB
	s3 *s3.Service
}

func NewProjectsRepository() *repository {
	return &repository{
		db: db.GetConnection(),
		s3: s3.NewS3Service(),
	}
}

/*
createProject creates a new project in the database
*/
func (r *repository) createProject(
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

// createProjectFiles creates the default files for a project in s3
func (r *repository) createProjectFiles(project model.Project) (model.ProjectFiles, error) {
	goMod := model.DefaultGoMod(project.Name)
	files := model.DefaultFiles

	files["go.mod"] = goMod

	srcDir := getProjectSrcDir(project.UserID, project.ID)

	if err := r.s3.UploadFiles(srcDir, files); err != nil {
		return nil, err
	}

	return model.DefaultFiles, nil
}

/*
getProjectsByUserID returns all projects for a given user (without files)
*/
func (r *repository) getProjectsByUserID(userID string) ([]model.ProjectView, error) {
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
func (r *repository) getProjectByID(userId string, id string) (model.ProjectView, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return model.ProjectView{}, err
	}

	if project.UserID != userId {
		return model.ProjectView{}, errors.New("project not found")
	}

	files, err := r.s3.GetFiles(getProjectSrcDir(userId, id))

	if err != nil {
		return model.ProjectView{}, err
	}

	wasmUrl, err := r.genProjectWasmPresignedURL(userId, id)
	if err != nil {
		return model.ProjectView{}, err
	}

	return project.ViewWith(wasmUrl, files), nil
}

/*
deleteProject deletes a project from the database
*/
func (r *repository) deleteProject(userId string, id string) error {

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
deleteProjectFiles deletes a project stored in s3
*/
func (r *repository) deleteProjectFiles(userId string, id string) error {
	dir := getProjectDir(userId, id)
	s3Err := r.s3.DeleteDir(dir)
	if s3Err != nil {
		return s3Err
	}

	return nil
}

/*
updateProjectFiles updates the files for a given project in s3
*/
func (r *repository) updateProjectSrcFiles(userId string, id string, files model.ProjectFiles) (
	model.ProjectFiles,
	error,
) {
	srcDir := getProjectSrcDir(userId, id)
	err := r.s3.UploadFiles(srcDir, files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *repository) uploadProjectWasm(userId string, id string, file io.Reader) error {
	wasmDir := getProjectWasmDir(userId, id)
	_, err := r.s3.Upload(wasmDir, "main.wasm", file)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) genProjectWasmPresignedURL(userId string, id string) (string, error) {
	wasmDir := getProjectWasmDir(userId, id)
	url, err := r.s3.GenPresignedURL(path.Join(wasmDir, "main.wasm"), time.Hour*24*7)
	if err != nil {
		return "", err
	}

	return url, nil
}

/* genSrcPresignedURL generates a presigned url for a source file in a project (not wasm) */
func (r *repository) genSrcPresignedURL(userId string, id string, filename string) (string, error) {
	srcDir := getProjectSrcDir(userId, id)
	url, err := r.s3.GenPresignedURL(path.Join(srcDir, filename), time.Hour*24*7)
	if err != nil {
		return "", err
	}

	return url, nil
}
