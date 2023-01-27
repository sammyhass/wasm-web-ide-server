package projects

import (
	"errors"
	"io"
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
	projUserId = "user_id"
)

var fakeProject = model.Project{
	ID:     "1",
	Name:   projName,
	UserID: projUserId,
	Model: &gorm.Model{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func (r *repoMock) createProject(name, userID string) (model.Project, error) {
	args := r.Called(name, userID)
	return fakeProject, args.Error(1)
}

func (r *repoMock) getProjectsByUserID(userId string) ([]model.ProjectView, error) {
	args := r.Called(userId)

	return []model.ProjectView{
		fakeProject.View(),
	}, args.Error(1)
}

func (r *repoMock) createProjectFiles(project model.Project) (model.ProjectFiles, error) {
	args := r.Called(project)
	return model.DefaultFiles, args.Error(1)
}

func (r *repoMock) getProjectByID(userId, projectID string) (model.ProjectView, error) {
	args := r.Called(userId, projectID)
	return fakeProject.ViewWithFiles(model.DefaultFiles), args.Error(1)
}

func (r *repoMock) deleteProject(userId, projectID string) error {
	args := r.Called(userId, projectID)
	return args.Error(0)
}

func (r *repoMock) deleteProjectFiles(userId, projectID string) error {
	args := r.Called(projectID)
	return args.Error(0)
}

func (r *repoMock) updateProjectSrcFiles(userId, projectId string, files model.ProjectFiles) (model.ProjectFiles, error) {
	args := r.Called(userId, projectId, files)
	return model.DefaultFiles, args.Error(1)
}

func (r *repoMock) uploadProjectWasm(userId, projectId string, f io.Reader) error {
	args := r.Called(userId, projectId, f)
	return args.Error(0)
}

func (r *repoMock) genPresignedURL(userId, projectId, fileName string) (string, error) {
	args := r.Called(userId, projectId, fileName)
	return "", args.Error(0)
}

func TestProjectsService_CreateProject_Success(t *testing.T) {
	repoMock := &repoMock{}

	s := &service{
		repo: repoMock,
	}

	repoMock.On("createProject", projName, projUserId).Return(fakeProject, nil)
	repoMock.On("createProjectFiles", fakeProject).Return(model.DefaultFiles, nil)

	pv, err := s.createProject(projName, projUserId)

	if err != nil {
		t.Error(err)
	}

	if pv.Name != projName {
		t.Error("Expected project name to be", projName)
	}

	if pv.UserID != projUserId {
		t.Error("Expected project user id to be", projUserId)
	}
	repoMock.AssertExpectations(t)
}

func TestProjectsService_CreateProject_DB_Error(t *testing.T) {
	repoMock := &repoMock{}
	dbErr := errors.New("Couldn't create project in DB")
	repoMock.On("createProject", projName, projUserId).Return(model.Project{}, dbErr)

	s := &service{
		repo: repoMock,
	}

	_, err := s.createProject("projectName", "user_id")

	if err != dbErr {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_CreateProject_S3_Error(t *testing.T) {
	repoMock := &repoMock{}
	s3Error := errors.New("Couldn't create project in S3")
	repoMock.On("createProject", projName, projUserId).Return(model.ProjectView{}, nil)
	repoMock.On("createProjectFiles", mock.Anything).Return(model.ProjectFiles{}, s3Error)

	s := &service{
		repo: repoMock,
	}

	_, err := s.createProject("projectName", "user_id")

	if err != s3Error {
		t.Error("Expected s3 error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectsByUserID_Success(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("getProjectsByUserID", projUserId).Return([]model.ProjectView{}, nil)

	s := &service{
		repo: repoMock,
	}

	pv, err := s.getProjectsByUserID("user_id")

	if err != nil {
		t.Error(err)
	}

	if pv[0].Name != "projectName" {
		t.Error("Expected project name to be projectName")
	}

	if pv[0].UserID != "user_id" {
		t.Error("Expected project user id to be user_id")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectsByUserID_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("getProjectsByUserID", projUserId).Return([]model.ProjectView{}, errors.New("Couldn't get projects"))

	s := &service{
		repo: repoMock,
	}

	_, err := s.getProjectsByUserID("user_id")

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_GetProjectByID(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("getProjectByID", projUserId, "1").Return(fakeProject.ViewWithFiles(model.DefaultFiles), nil)

	s := &service{
		repo: repoMock,
	}

	_, err := s.getProjectByID(projUserId, "1")

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("deleteProject", projUserId, "1").Return(nil)
	repoMock.On("deleteProjectFiles", "1").Return(nil)

	s := &service{
		repo: repoMock,
	}

	err := s.deleteProjectByID(projUserId, "1")

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID_DB_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("deleteProject", projUserId, "1").Return(errors.New("Couldn't delete project"))

	s := &service{
		repo: repoMock,
	}

	err := s.deleteProjectByID(projUserId, "1")

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_DeleteProjectByID_S3_Error(t *testing.T) {
	s3Err := errors.New("Couldn't delete project files")
	repoMock := &repoMock{}
	repoMock.On("deleteProject", projUserId, "1").Return(nil)
	repoMock.On("deleteProjectFiles", "1").Return(s3Err)

	s := &service{
		repo: repoMock,
	}

	err := s.deleteProjectByID(projUserId, "1")

	if err != s3Err {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_UpdateProjectFiles_Success(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("updateProjectSrcFiles", projUserId, "1", model.DefaultFiles).Return(model.DefaultFiles, nil)

	s := &service{
		repo: repoMock,
	}

	_, err := s.updateProjectFiles(projUserId, "1", model.DefaultFiles)

	if err != nil {
		t.Error(err)
	}

	repoMock.AssertExpectations(t)
}

func TestProjectsService_UpdateProjectFiles_Error(t *testing.T) {
	repoMock := &repoMock{}
	repoMock.On("updateProjectSrcFiles", projUserId, "1", model.DefaultFiles).Return(model.DefaultFiles, errors.New("Couldn't update project files"))

	s := &service{
		repo: repoMock,
	}

	_, err := s.updateProjectFiles(projUserId, "1", model.DefaultFiles)

	if err == nil {
		t.Error("Expected error to be returned")
	}

	repoMock.AssertExpectations(t)
}
