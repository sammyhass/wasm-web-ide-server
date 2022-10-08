package main

import (
	"github.com/sammyhass/web-ide/server/modules/env"
	"github.com/sammyhass/web-ide/server/modules/router"
)

func main() {
	env.InitEnv()

	router.Run(env.Env.PORT)
}
