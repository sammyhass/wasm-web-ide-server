package router

import (
	"github.com/sammyhass/web-ide/server/modules/auth"
	"github.com/sammyhass/web-ide/server/modules/file_server"
	"github.com/sammyhass/web-ide/server/modules/projects"
)

func Run(
	port string,
) {
	router := NewRouter()

	router.UseController(file_server.CONTROLLER_ROUTE, file_server.NewController())
	router.UseController("/auth", auth.NewController())
	router.UseController("/projects", projects.NewController())

	router.Routes()

	router.Run(port)
}
