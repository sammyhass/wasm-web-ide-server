package router

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sammyhass/web-ide/server/env"
	"github.com/sammyhass/web-ide/server/router/middleware"
)

// controller should be implemented by all controllers in order to register their routes with the gin engine
type controller interface {
	// Routes can be used to register routes for a given controller
	Routes(e *gin.RouterGroup)
}

// router is a wrapper around the gin.Engine that allows for the registration of RouterGroup
type router struct {
	Engine *gin.Engine

	controllers map[string]controller
}

func newEngine() *gin.Engine {
	eng := gin.Default()
	eng.RedirectTrailingSlash = false
	eng.RedirectFixedPath = false

	return eng
}

func newRouter() *router {
	return &router{
		Engine:      newEngine(),
		controllers: make(map[string]controller),
	}
}

// useController register a controller with the router
func (r *router) useController(name string, controller controller) {
	if _, ok := r.controllers[name]; ok {
		log.Fatalf("controller %s already registered", name)

	}

	r.controllers[name] = controller
}

// Run starts the server on the given port
func (r *router) run(port string) {
	r.Engine.Run(
		":" + port,
	)
}

// routes runs the routes function for each controller with a router group
func (r *router) routes() {
	r.Middleware()

	for name, controller := range r.controllers {
		group := r.Engine.Group(name)
		controller.Routes(group)
	}
}

// Middleware should be used to register all middleware for the router
func (r *router) Middleware() {

	r.useCORS()

	r.Engine.Use(middleware.ErrorHandler)
	r.Engine.Use(middleware.AuthMiddleware)
}

func (r *router) useCORS() {
	allowedHeaders := []string{"Origin", "Content-Length", "Content-Type", "Authorization", "User-Agent", "Referer", "Cache-Control", "X-Requested-With",
		"Access-Control-Request-Headers", "Access-Control-Request-Method", "Accept-Encoding", "Accept-Language", "Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site", "Sec-Fetch-User", "Host", "Connection", "Upgrade-Insecure-Requests", "Cache-Control", "Accept", "Accept-Encoding", "Accept-Language", "User-Agent", "Pragma"}

	corsOrigin := env.GetOr(env.CORS_ALLOW_ORIGIN, "http://localho.st:3000")

	fmt.Println("CORS origin:", corsOrigin)

	r.Engine.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{corsOrigin},
				AllowCredentials: true,
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
				AllowHeaders:     allowedHeaders,
			},
		),
	)
}
