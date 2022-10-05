package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/db"
)

type ProjectsController struct {
	service ProjectsService
}

func NewController() ProjectsController {
	svc := NewService(db.GetDB())

	return ProjectsController{
		service: svc,
	}
}

func (pr ProjectsController) Routes(e *gin.Engine) {
	e.GET("/projects", pr.getProjects)
	e.POST("/projects", pr.createProject)
}

func (pr ProjectsController) getProjects(c *gin.Context) {
	pr.service.GetProjects()
}

type createProjectDTO struct {
	Name string `json:"name"`
}

func (pr ProjectsController) createProject(c *gin.Context) {
	var dto createProjectDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err)
		return
	}

	proj, err := pr.service.CreateProject(dto.Name)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, proj)
}
