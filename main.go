package main

import (
	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/env"
	"github.com/sammyhass/web-ide/server/modules/migrations"
	"github.com/sammyhass/web-ide/server/modules/router"
)

func main() {
	env.InitEnv()
	db.Connect()

	migrations.Migrate()

	router.Run(env.Env.PORT)
}
