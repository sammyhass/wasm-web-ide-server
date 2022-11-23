package main

import (
	"github.com/sammyhass/web-ide/server/modules/db"
	"github.com/sammyhass/web-ide/server/modules/env"
	"github.com/sammyhass/web-ide/server/modules/model"
	"github.com/sammyhass/web-ide/server/modules/router"
	"github.com/sammyhass/web-ide/server/modules/s3"
)

func main() {
	env.InitEnv()

	db.Connect()
	defer db.Close()

	s3.InitSession()

	model.Migrate()

	router.Run(env.Get(env.PORT))
}
