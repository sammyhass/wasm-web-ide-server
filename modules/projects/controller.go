package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/auth"
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
	group.GET("/", auth.Protected(c.getProjects))
	group.POST("/", auth.Protected(c.createProject))

	group.GET("/:id", auth.Protected(c.getProject))
	group.POST("/:id/compile")
}

type newProjectDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *ProjectsController) createProject(
	ctx *gin.Context,
	uuid string,
) {
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

	ctx.JSON(200, proj)
}

func (c *ProjectsController) getProjects(
	ctx *gin.Context,
	uuid string,
) {
	projects, err := c.service.GetProjectsByUserID(uuid)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, projects)
}

func (c *ProjectsController) getProject(
	ctx *gin.Context,
	uuid string,
) {
	project, err := c.service.GetProjectByID(uuid, ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, project)
}
