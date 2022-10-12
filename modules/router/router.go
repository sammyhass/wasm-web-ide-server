package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sammyhass/web-ide/server/modules/router/middleware"
)

// Controller should be implemented by all controllers in order to register their routes with the gin engine
type Controller interface {

	// Routes allows the group to register its routes with the router
	Routes(e *gin.Engine)
}

// Router is a wrapper around the gin.Engine that allows for the registration of RouterGroup
type Router struct {
	Engine *gin.Engine

	controllers []Controller
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}

// UseController adds a controller to the router which will be registered when Routes is called
func (r *Router) UseController(controller Controller) {
	r.controllers = append(r.controllers, controller)
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
	for _, group := range r.controllers {
		group.Routes(r.Engine)
	}
}

// Middleware should be used to register all middleware for the router
func (r Router) Middleware() {
	r.Engine.Use(cors.Default())
	r.Engine.Use(middleware.ErrorHandler)
}
