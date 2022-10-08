package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sammyhass/web-ide/server/modules/router/middleware"
)

// RouterGroup is an interface that should be implemented by all groups that are added to the router.
type RouterGroup interface {

	// Routes allows the group to register its routes with the router
	Routes(e *gin.Engine)
}

// Router is a wrapper around the gin.Engine that allows for the registration of RouterGroup
type Router struct {
	Engine *gin.Engine

	groups []RouterGroup
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}

// AddGroup adds a RouterGroup to the router which will be registered when Routes is called
func (r *Router) AddGroup(group RouterGroup) {
	r.groups = append(r.groups, group)
}

// Run starts the server on the given port
func (r *Router) Run(port string) {
	r.Engine.Run(
		":" + port,
	)
}

// Routes runs the Routes function for each group that has been registered
func (r *Router) Routes() {
	r.Middleware()
	for _, group := range r.groups {
		group.Routes(r.Engine)
	}
}

// Middleware should be used to register all middleware for the router
func (r Router) Middleware() {
	r.Engine.Use(cors.Default())
	r.Engine.Use(middleware.ErrorHandler)
}
