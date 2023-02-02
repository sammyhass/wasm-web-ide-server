package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/auth"
	"github.com/sammyhass/web-ide/server/modules/model"
)

type controller struct {
	service *service
}

func NewController() *controller {
	return &controller{
		service: newService(),
	}
}

func (c *controller) Routes(
	group *gin.RouterGroup,
) {
	group.GET("", auth.Protected(c.getProjects))
	group.POST("", auth.Protected(c.createProject))

	group.GET("/:id", auth.Protected(c.getProject))
	group.DELETE("/:id", auth.Protected(c.deleteProject))
	group.PATCH("/:id", auth.Protected(c.updateProject))
	group.POST("/:id/compile", auth.Protected(c.compileProjectToWasm))
	group.GET("/:id/wat", auth.Protected(c.getProjectWat))

}

type newProjectDto struct {
	Name string `json:"name"`
}

func (c *controller) createProject(
	ctx *gin.Context,
	uuid string,
) {
	var dto newProjectDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	proj, err := c.service.createProject(dto.Name, uuid)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, proj)
}

func (c *controller) getProjects(
	ctx *gin.Context,
	uuid string,
) {
	projects, err := c.service.getProjectsByUserID(uuid)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, projects)
}

func (c *controller) getProject(
	ctx *gin.Context,
	uuid string,
) {
	project, err := c.service.getProjectByID(uuid, ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, project)
}

func (c *controller) deleteProject(
	ctx *gin.Context,
	uuid string,
) {
	err := c.service.deleteProjectByID(uuid, ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{})

}

type updateProjectFilesDto struct {
	Files model.ProjectFiles `json:"files"`
}

func (c *controller) updateProject(
	ctx *gin.Context,
	uuid string,
) {
	var dto updateProjectFilesDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	_, err := c.service.updateProjectFiles(
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

func (c *controller) compileProjectToWasm(
	ctx *gin.Context,
	uuid string,
) {
	path, err := c.service.compileProjectWASM(uuid, ctx.Param("id"))

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, path)
}

func (c *controller) getProjectWat(
	ctx *gin.Context,
	uuid string,
) {
	wat, err := c.service.genProjectWatPresignedURL(uuid, ctx.Param("id"))

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, wat)
}
