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
	group.GET("/", auth.Protected(c.getProjects))
	group.POST("/", auth.Protected(c.createProject))

	group.GET("/:id", auth.Protected(c.getProject))
	group.DELETE("/:id", auth.Protected(c.deleteProject))
	group.PATCH("/:id", auth.Protected(c.updateProject))
	group.POST("/:id/compile", auth.Protected(c.compileProjectToWasm))

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

func (c *ProjectsController) deleteProject(
	ctx *gin.Context,
	uuid string,
) {
	err := c.service.DeleteProjectByID(uuid, ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(200)

}

type updateProjectFilesDto struct {
	Files model.ProjectFiles `json:"files"`
}

func (c *ProjectsController) updateProject(
	ctx *gin.Context,
	uuid string,
) {
	var dto updateProjectFilesDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	_, err := c.service.UpdateProjectFiles(
		uuid,
		ctx.Param("id"),
		dto.Files,
	)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200,
		model.ProjectFilesToFileViews(dto.Files),
	)
}

func (c *ProjectsController) compileProjectToWasm(
	ctx *gin.Context,
	uuid string,
) {
	path, err := c.service.CompileProjectWASM(uuid, ctx.Param("id"))

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, path)
}
