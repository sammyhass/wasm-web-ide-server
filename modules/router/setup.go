package router

import (
	"github.com/sammyhass/web-ide/server/modules/auth"
	"github.com/sammyhass/web-ide/server/modules/file_server"
	"github.com/sammyhass/web-ide/server/modules/wasm"
)

func Run(
	port string,
) {
	router := NewRouter()

	router.UseController(wasm.CONTROLLER_ROUTE, wasm.NewController())
	router.UseController(file_server.CONTROLLER_ROUTE, file_server.NewController())
	router.UseController("/auth", auth.NewController())

	router.Routes()

	router.Run(port)
}
