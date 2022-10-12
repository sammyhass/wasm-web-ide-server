package router

import (
	"github.com/sammyhass/web-ide/server/modules/file_server"
	"github.com/sammyhass/web-ide/server/modules/wasm"
)

func Run(
	port string,
) {
	router := NewRouter()

	router.Use(wasm.NewController())
	router.Use(file_server.NewController())

	router.Routes()

	router.Run(port)
}
