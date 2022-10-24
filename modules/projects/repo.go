package projects

import (
	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/model"
	"gorm.io/gorm"
)

type ProjectsRepository struct {
	db *gorm.DB
}

func NewProjectsRepository() *ProjectsRepository {
	return &ProjectsRepository{
		db: db.GetConnection(),
	}
}

func (r *ProjectsRepository) CreateProject(
	name string,
	description string,
	userID string,
) (model.Project, error) {

	proj := model.Project{
		ID:          model.NewID(),
		Name:        name,
		Description: description,
		UserID:      userID,
	}

	err := r.db.Create(&proj).Error

	return proj, err

}

func (r *ProjectsRepository) GetProjectsByUserID(userID string) ([]*model.Project, error) {
	var projects []*model.Project

	err := r.db.Where("user_id = ?", userID).Find(&projects).Error

	return projects, err
}

func (r *ProjectsRepository) GetProjectByID(id string) (*model.Project, error) {
	var project model.Project

	err := r.db.Where("id = ?", id).First(&project).Error

	return &project, err
}
