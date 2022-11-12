package projects

import (
	"errors"
	"testing"
	"time"

	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type repoMock struct {
	mock.Mock
	projectsRepo
}

const (
	// ProjectName is the name of the project
	projName   = "projectName"
	projDesc   = "test"
	projUserId = "user_id"
)

var fakeProject = model.Project{
	ID:          "1",
	Name:        projName,
	Description: projDesc,
	UserID:      projUserId,
	Model: &gorm.Model{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func (r *repoMock) CreateProject(name, description, userID string) (model.Project, error) {
	args := r.Called(name, description, userID)
	return fakeProject, args.Error(1)
}

func (r *repoMock) GetProjectsByUserID(userId string) ([]model.ProjectView, error) {
	args := r.Called(userId)

	return []model.ProjectView{
		fakeProject.View(),
	}, args.Error(1)
}

func (r *repoMock) CreateProjectFiles(project model.Project) (model.ProjectFiles, error) {
	args := r.Called(project)
	return model.DefaultFiles, args.Error(1)
}

func (r *repoMock) GetProjectByID(userId, projectID string) (model.ProjectView, error) {
	args := r.Called(userId, projectID)
	return fakeProject.ViewWithFiles(model.DefaultFiles), args.Error(1)
}

func (r *repoMock) DeleteProject(userId, projectID string) error {
	args := r.Called(userId, projectID)
	return args.Error(0)
}

func (r *repoMock) DeleteProjectFiles(userId, projectID string) error {
	args := r.Called(projectID)
	return args.Error(0)
}

func (r *repoMock) UpdateProjectFiles(userId, projectId string, files model.ProjectFiles) (model.ProjectFiles, error) {
	args := r.Called(userId, projectId, files)
	return model.DefaultFiles, args.Error(1)
}

func TestProjectsService_CreateProject_Success(t *testing.T) {
	repoMock := &repoMock{}

	s := &ProjectsService{
		repo: repoMock,
	}

	repoMock.On("CreateProject", projName, projDesc, projUserId).Return(fakeProject, nil)
	repoMock.On("CreateProjectFiles", fakeProject).Return(model.DefaultFiles, nil)

	pv, err := s.CreateProject(projName, projDesc, projUserId)

	if err != nil {
		t.Error(err)
	}

	if pv.Name != projName {
		t.Error("Expected project name to be", projName)
	}

	if pv.Description != projDesc {
		t.Error("Expected project description to be", projDesc)
	}

	if pv.UserID != projUserId {
		t.Error("Expected project user id to be", projUserId)
	}
	repoMock.AssertExpectations(t)
}

func TestProjectsService_CreateProject_DB_Error(t *testing.T) {
	repoMock := &repoMock{}
	dbErr := errors.New("Couldn't create project in DB")
	repoMock.On("CreateProject", projName, projDesc, projUserId).Return(model.Project{}, dbErr)

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.CreateProject("projectName", "test", "user_id")

	if err != dbErr {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_CreateProject_S3_Error(t *testing.T) {
	repoMock := &repoMock{}
	s3Error := errors.New("Couldn't create project in S3")
	repoMock.On("CreateProject", projName, projDesc, projUserId).Return(model.ProjectView{}, nil)
	repoMock.On("CreateProjectFiles", mock.Anything).Return(model.ProjectFiles{}, s3Error)

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.CreateProject("projectName", "test", "user_id")

	if err != s3Error {
		t.Error("Expected s3 error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectsByUserID_Success(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("GetProjectsByUserID", projUserId).Return([]model.ProjectView{}, nil)

	s := &ProjectsService{
		repo: repoMock,
	}

	pv, err := s.GetProjectsByUserID("user_id")

	if err != nil {
		t.Error(err)
	}

	if pv[0].Name != "projectName" {
		t.Error("Expected project name to be projectName")
	}

	if pv[0].Description != "test" {
		t.Error("Expected project description to be test")
	}

	if pv[0].UserID != "user_id" {
		t.Error("Expected project user id to be user_id")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectsByUserID_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("GetProjectsByUserID", projUserId).Return([]model.ProjectView{}, errors.New("Couldn't get projects"))

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.GetProjectsByUserID("user_id")

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectByID(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("GetProjectByID", projUserId, "1").Return(fakeProject.ViewWithFiles(model.DefaultFiles), nil)

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.GetProjectByID(projUserId, "1")

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("DeleteProject", projUserId, "1").Return(nil)
	repoMock.On("DeleteProjectFiles", "1").Return(nil)

	s := &ProjectsService{
		repo: repoMock,
	}

	err := s.DeleteProjectByID(projUserId, "1")

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID_DB_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("DeleteProject", projUserId, "1").Return(errors.New("Couldn't delete project"))

	s := &ProjectsService{
		repo: repoMock,
	}

	err := s.DeleteProjectByID(projUserId, "1")

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID_S3_Error(t *testing.T) {
	s3Err := errors.New("Couldn't delete project files")
	repoMock := &repoMock{}
	repoMock.On("DeleteProject", projUserId, "1").Return(nil)
	repoMock.On("DeleteProjectFiles", "1").Return(s3Err)

	s := &ProjectsService{
		repo: repoMock,
	}

	err := s.DeleteProjectByID(projUserId, "1")

	if err != s3Err {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_UpdateProjectFiles_Success(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("UpdateProjectFiles", projUserId, "1", model.DefaultFiles).Return(model.DefaultFiles, nil)

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.UpdateProjectFiles(projUserId, "1", model.DefaultFiles)

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_UpdateProjectFiles_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("UpdateProjectFiles", projUserId, "1", model.DefaultFiles).Return(model.DefaultFiles, errors.New("Couldn't update project files"))

	s := &ProjectsService{
		repo: repoMock,
	}

	_, err := s.UpdateProjectFiles(projUserId, "1", model.DefaultFiles)

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}
