package file_server

import (
	"os"

	"github.com/gin-gonic/gin"
)

const (
	/*
		STATIC_DIR is the directory where static files are stored on the server
	*/
	STATIC_DIR string = "./www"

	/*
		STATIC_ROUTE is the route where static files are served from
	*/
	STATIC_ROUTE string = "/static"
)

type StaticFileServerController struct{}

func NewController() *StaticFileServerController {
	return &StaticFileServerController{}
}

func (ssc *StaticFileServerController) Routes(e *gin.Engine) {
	// Ensure the directory we are serving from exists.
	if _, err := os.Stat(STATIC_DIR); os.IsNotExist(err) {
		os.Mkdir(STATIC_DIR, 0755)
	}

	e.Static(STATIC_ROUTE, STATIC_DIR)
}
