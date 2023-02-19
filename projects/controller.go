package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/auth"
	"github.com/sammyhass/web-ide/server/model"
)

type controller struct {
	service *Service
}

func NewController() *controller {
	return &controller{
		service: NewService(),
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
	group.PATCH("/:id/rename", auth.Protected(c.renameProject))
	group.PATCH("/:id/share", auth.Protected(c.toggleShareProject))

	group.POST("/fork/:code", auth.Protected(c.forkProject))
	group.GET("/fork/:code", c.getSharedProject)
}

type newProjectDto struct {
	Name     string `json:"name"`
	Language string `json:"language"`
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

	lang := model.GetProjectLanguage(dto.Language)

	proj, err := c.service.CreateProject(dto.Name, uuid, lang)

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
	projects, err := c.service.GetProjectsByUserID(uuid)
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
	project, err := c.service.GetProjectByID(uuid, ctx.Param("id"))
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
	err := c.service.DeleteProjectByID(uuid, ctx.Param("id"))
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

func (c *controller) compileProjectToWasm(
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

func (c *controller) getProjectWat(
	ctx *gin.Context,
	uuid string,
) {
	wat, err := c.service.GenProjectWatPresignedURL(uuid, ctx.Param("id"))

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, wat)
}

type renameProjectDto struct {
	Name string `json:"name"`
}

func (c *controller) renameProject(
	ctx *gin.Context,
	uuid string,
) {
	id := ctx.Param("id")

	var dto renameProjectDto

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	if p, err := c.service.RenameProject(uuid, id, dto.Name); err != nil {
		ctx.Error(err)
		return
	} else {
		ctx.JSON(200, p)
	}
}

// create a share code by which a project can be forked by another user
// returns the share code if the project is now shareable, else returns false
func (c *controller) toggleShareProject(
	ctx *gin.Context,
	uuid string,
) {

	id := ctx.Param("id")

	var body struct {
		Shared bool `json:"shared"`
	} = struct {
		Shared bool `json:"shared"`
	}{}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(err)
		return
	}

	shareCode, err := c.service.ToggleForkable(uuid, id, body.Shared)

	if err != nil {
		ctx.Error(err)
		return
	}

	res := make(map[string]interface{})
	if shareCode == "" {
		res["shared"] = false
	} else {
		res["shared"] = true
		res["shareCode"] = shareCode
	}

	ctx.JSON(200, res)
}

func (c *controller) forkProject(
	ctx *gin.Context,
	uuid string,
) {
	shareCode := ctx.Param("code")

	proj, err := c.service.ForkProject(uuid, shareCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, proj)
}

func (c *controller) getSharedProject(
	ctx *gin.Context,
) {
	shareCode := ctx.Param("code")

	project, err := c.service.GetSharedProject(shareCode)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, project)
}
