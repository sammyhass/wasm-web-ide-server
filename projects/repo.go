package projects

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/sammyhass/web-ide/server/db"
	"github.com/sammyhass/web-ide/server/model"
	"github.com/sammyhass/web-ide/server/s3"
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

type Repository struct {
	db *gorm.DB
	s3 *s3.Service
}

func newRepository() *Repository {
	return &Repository{
		db: db.GetConnection(),
		s3: s3.NewService(),
	}
}

/*
createProject creates a new project in the database
*/
func (r *Repository) createProject(
	name string,
	userID string,
	language model.ProjectLanguage,
) (model.Project, error) {

	proj := model.NewProject(
		name,
		userID,
		language,
	)

	err := r.db.Create(&proj).Error

	if err != nil {
		return model.Project{}, err
	}

	return proj, nil
}

// createProjectFilesWith creates a project in the s3 bucket with the given files
func (r *Repository) createProjectFilesWith(userId, projectId string, files model.ProjectFiles) (model.ProjectFiles, error) {
	srcDir := getProjectSrcDir(userId, projectId)

	if files == nil {
		return nil, errors.New("invalid project language")
	}

	if err := r.s3.UploadFiles(srcDir, files); err != nil {
		return nil, err
	}

	return files, nil
}

// createProjectFiles creates a project in the s3 bucket with the default files for the given language
func (r *Repository) createProjectFiles(userId, projectId string, language model.ProjectLanguage) (model.ProjectFiles, error) {
	var files model.ProjectFiles
	switch language {
	case model.LanguageGo:
		files = model.DefaultFilesGo
	case model.LanguageAssemblyScript:
		files = model.DefaultFilesAssemblyScript
	default:
		return nil, errors.New("invalid project language")
	}

	r.createProjectFilesWith(userId, projectId, files)

	return files, nil
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

func (r *Repository) getProjectRecord(
	userId string,
	id string,
) (model.Project, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return model.Project{}, err
	}

	if project.UserID != userId {
		return model.Project{}, errors.New("project not found")
	}

	return project, nil
}

/*
getProjectByID returns a project for a given user with the given id returning the view of the database record and the files in s3
*/
func (r *Repository) getProjectByID(userId string, id string) (model.ProjectView, error) {

	project, err := r.getProjectRecord(userId, id)

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

	return project.ViewWithFiles(files), nil
}

/*
deleteProject deletes a project from the database
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
deleteProjectFiles deletes a project stored in s3
*/
func (r *Repository) deleteProjectFiles(userId string, id string) error {
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
func (r *Repository) uploadProjectSrcFiles(userId string, id string, files model.ProjectFiles) (
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

func (r *Repository) uploadBuildFile(userId string, id string, name string, file io.Reader) error {
	wasmDir := getProjectWasmDir(userId, id)
	_, err := r.s3.Upload(wasmDir, name, file)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) uploadProjectWasm(userId string, id string, file io.Reader) error {
	return r.uploadBuildFile(userId, id, "main.wasm", file)
}

func (r *Repository) genProjectWasmPresignedURL(userId string, id string) (string, error) {
	wasmDir := getProjectWasmDir(userId, id)
	url, err := r.s3.GenPresignedURL(path.Join(wasmDir, "main.wasm"), time.Hour*24*7)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *Repository) uploadProjectWat(userId string, id string, file io.Reader) error {
	return r.uploadBuildFile(userId, id, "main.wat", file)
}

func (r *Repository) genProjectWatPresignedURL(userId string, id string) (string, error) {
	wasmDir := getProjectWasmDir(userId, id)
	url, err := r.s3.GenPresignedURL(path.Join(wasmDir, "main.wat"), time.Hour*24*7)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *Repository) renameProject(
	userId, id, name string,
) (model.ProjectView, error) {

	p, err := r.getProjectRecord(userId, id)
	if err != nil {
		return model.ProjectView{}, err
	}

	p.Name = name

	if err := r.db.Save(&p).Error; err != nil {
		return model.ProjectView{}, err
	}

	return p.View(), nil
}

func (r *Repository) allowSharing(p *model.Project) (sharecode string, err error) {
	p.IsShared = true
	p.ShareCode = sql.NullString{
		String: generateShareCode(),
		Valid:  true,
	}
	for {
		if err := r.db.Save(p).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				p.ShareCode = sql.NullString{
					String: generateShareCode(),
					Valid:  true,
				}
				continue
			}
			return "", err
		}
		break
	}

	return p.ShareCode.String, nil
}

func (r *Repository) toggleSharing(userId, id string, share bool) (sharecode string, err error) {
	p, err := r.getProjectRecord(userId, id)
	if err != nil {
		return "", err
	}

	if share {
		return r.allowSharing(&p)
	}

	return "", r.disallowSharing(&p)
}

func (r *Repository) disallowSharing(
	p *model.Project,
) error {

	if err := r.db.Model(p).UpdateColumns(
		map[string]interface{}{
			"is_shared":  false,
			"share_code": gorm.Expr("NULL"),
		},
	).Error; err != nil {
		return err
	}

	return nil
}

// getProjectByShareCode returns the project with the given share code files included
func (r *Repository) getProjectByShareCode(code string) (model.ProjectView, error) {
	var p model.Project
	if err := r.db.Where("share_code = ? AND is_shared = ?", code, true).First(&p).Error; err != nil {
		fmt.Println("getProjectByShareCodeError", err)
		return model.ProjectView{}, err
	}

	files, err := r.s3.GetFiles(getProjectSrcDir(p.UserID, p.ID))
	if err != nil {
		return p.View(), err
	}

	return p.ViewWithFiles(files), nil
}
