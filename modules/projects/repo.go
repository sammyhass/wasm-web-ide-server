package projects

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/s3"
	"gorm.io/gorm"
)

type projectsRepo interface {
	CreateProject(
		name string,
		userID string,
	) (model.Project, error)

	CreateProjectFiles(project model.Project) (model.ProjectFiles, error)

	GetProjectsByUserID(userID string) ([]model.ProjectView, error)
	GetProjectByID(userId string, id string) (model.ProjectView, error)

	DeleteProject(userId string, id string) error
	DeleteProjectFiles(userId string, id string) error

	UpdateProjectFiles(userId string, id string, files model.ProjectFiles) (model.ProjectFiles, error)
}

type ProjectsRepository struct {
	db *gorm.DB
	s3 *s3.S3Service
}

func NewProjectsRepository() *ProjectsRepository {
	return &ProjectsRepository{
		db: db.GetConnection(),
		s3: s3.NewS3Service(),
	}
}

/*
CreateProject creates a new project in the database
*/
func (r *ProjectsRepository) CreateProject(
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

/*
CreateProjectInS3 creates a new project in S3 for a given project stored in the database
*/
func (r *ProjectsRepository) CreateProjectFiles(project model.Project) (model.ProjectFiles, error) {
	_, err := r.s3.UploadProjectFiles(project.UserID, project.ID, model.DefaultFiles)
	if err != nil {
		return nil, err
	}

	return model.DefaultFiles, nil
}

/*
GetProjectsByUserID returns all projects for a given user (without files)
*/
func (r *ProjectsRepository) GetProjectsByUserID(userID string) ([]model.ProjectView, error) {
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
GetProjectByID returns a project for a given user with the given id returning the database record and the files in s3
*/
func (r *ProjectsRepository) GetProjectByID(userId string, id string) (model.ProjectView, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return model.ProjectView{}, err
	}

	if project.UserID != userId {
		return model.ProjectView{}, errors.New("project not found")
	}

	files, err := r.s3.GetProjectFiles(project.UserID, project.ID)

	if err != nil {
		return model.ProjectView{}, err
	}

	return project.ViewWithFiles(files), nil
}

/*
DeleteProjectDB deletes a project from the database
*/
func (r *ProjectsRepository) DeleteProject(userId string, id string) error {

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
DeleteProjectFiles deletes a project stored in s3
*/
func (r *ProjectsRepository) DeleteProjectFiles(userId string, id string) error {
	s3Err := r.s3.DeleteProjectFiles(userId, id)
	if s3Err != nil {
		return s3Err
	}

	return nil
}

/*
UpdateProjectFiles updates the files for a given project in s3
*/
func (r *ProjectsRepository) UpdateProjectFiles(userId string, id string, files model.ProjectFiles) (
	model.ProjectFiles,
	error,
) {
	_, err := r.s3.UploadProjectFiles(userId, id, files)
	if err != nil {
		return nil, err
	}

	return files, nil
}
