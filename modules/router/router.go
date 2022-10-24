package router

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sammyhass/web-ide/server/modules/router/middleware"
)

// Controller should be implemented by all controllers in order to register their routes with the gin engine
type Controller interface {
	// Routes can be used to register routes for a given controller
	Routes(e *gin.RouterGroup)
}

// Router is a wrapper around the gin.Engine that allows for the registration of RouterGroup
type Router struct {
	Engine *gin.Engine

	controllers map[string]Controller
}

func NewRouter() *Router {
	return &Router{
		Engine:      gin.Default(),
		controllers: make(map[string]Controller),
	}
}

// UseController register a controller with the router
func (r *Router) UseController(name string, controller Controller) {
	if _, ok := r.controllers[name]; ok {
		log.Fatalf("controller %s already registered", name)
	}

	r.controllers[name] = controller
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

	for name, controller := range r.controllers {
		group := r.Engine.Group(name)
		controller.Routes(group)
	}
}

// Middleware should be used to register all middleware for the router
func (r Router) Middleware() {
	r.Engine.Use(
		cors.New(
			cors.Config{
				AllowOrigins: []string{"http://localho.st:3000"},
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
				AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			},
		),
	)

	r.Engine.Use(
		func(ctx *gin.Context) {
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(200)
				return
			}

			ctx.Next()
		},
	)
	r.Engine.Use(middleware.ErrorHandler)
}
