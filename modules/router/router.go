package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sammyhass/web-ide/server/modules/router/middleware"
)

type RouterGroup interface {
	// Routes registers the routes for the group
	Routes(e *gin.Engine)
}

type Router struct {
	Engine *gin.Engine

	groups []RouterGroup
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}

func (r *Router) AddGroup(group RouterGroup) {
	r.groups = append(r.groups, group)
}

func (r *Router) Run(port string) {
	r.Engine.Run(
		":" + port,
	)
}

func (r *Router) Routes() {
	r.Middleware()
	for _, group := range r.groups {
		group.Routes(r.Engine)
	}
}

func (r Router) Middleware() {
	r.Engine.Use(cors.Default())
	r.Engine.Use(middleware.ErrorHandler)
}
