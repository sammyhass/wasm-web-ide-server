package projects

import (
	"errors"

	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/s3"
	"gorm.io/gorm"
)

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

func (r *ProjectsRepository) CreateProject(
	name string,
	description string,
	userID string,
) (model.ProjectView, error) {

	proj := model.NewProject(
		name,
		description,
		userID,
	)

	err := r.db.Create(&proj).Error

	if err != nil {
		return model.ProjectView{}, err
	}

	files, err := r.createProjectInS3(proj)
	if err != nil {
		return model.ProjectView{}, err
	}

	fviews := model.ProjectFilesToFileViews(files)

	projView := proj.ViewWithFiles(fviews)

	return projView, nil

}

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

func (r *ProjectsRepository) GetProjectByID(userId string, id string) (model.ProjectView, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	if err != nil {
		return model.ProjectView{}, err
	}

	if project.UserID != userId {
		return model.ProjectView{}, errors.New("project not found")
	}

	files, err := r.s3.GetFiles(project.UserID, project.ID)

	if err != nil {
		return model.ProjectView{}, err
	}

	return project.ViewWithFiles(
		model.ProjectFilesToFileViews(files),
	), nil
}

/*
createProjectInS3 initializes a users project in S3, creating app.js, main.go, index.html, styles.css
When saved, these files are stored at `[userId]/[projectId]/[fileName]`
*/
func (r *ProjectsRepository) createProjectInS3(project model.Project) (
	model.ProjectFiles,
	error,
) {
	_, err := r.s3.UploadFiles(project.UserID, project.ID, model.DefaultFiles)

	if err != nil {
		return nil, err
	}

	return model.DefaultFiles, nil

}
