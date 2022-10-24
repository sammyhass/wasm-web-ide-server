package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/auth"
	"github.com/sammyhass/web-ide/server/modules/model"
)

type ProjectsController struct {
	service *ProjectsService
}

func NewController() *ProjectsController {
	return &ProjectsController{
		service: NewProjectsService(),
	}
}

func (c *ProjectsController) Routes(
	group *gin.RouterGroup,
) {
	group.POST("/", c.createProject)
	group.GET("/:id", c.getProject)
	group.GET("/", c.getProjects)
}

type newProjectDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *ProjectsController) createProject(
	ctx *gin.Context,
) {
	uuid := auth.GetUserFromContextOrAbort(ctx)

	var dto newProjectDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	proj, err := c.service.CreateProject(dto.Name, dto.Description, uuid)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, proj.View())
}

func (c *ProjectsController) getProjects(
	ctx *gin.Context,
) {
	uuid := auth.GetUserFromContextOrAbort(ctx)

	projects, err := c.service.GetProjectsByUserID(uuid)

	if err != nil {
		ctx.Error(err)
		return
	}

	views := make([]model.ProjectView, len(projects))
	for i, p := range projects {
		views[i] = p.View()
	}

	ctx.JSON(200, projects)
}

func (c *ProjectsController) getProject(
	ctx *gin.Context,
) {
	uuid := auth.GetUserFromContextOrAbort(ctx)

	project, err := c.service.GetProjectByID(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	if project.UserID != uuid {
		ctx.AbortWithStatus(403)
		return
	}

	ctx.JSON(200, project.View())
}
