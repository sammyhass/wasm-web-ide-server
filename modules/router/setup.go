package router

import (
	"github.com/sammyhass/web-ide/server/modules/auth"
	"github.com/sammyhass/web-ide/server/modules/projects"
)

func Run(
	port string,
) {
	router := newRouter()

	router.useController("/auth", auth.NewController())
	router.useController("/projects", projects.NewController())

	router.routes()

	router.run(port)
}
