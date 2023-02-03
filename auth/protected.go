package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Protected(handler func(ctx *gin.Context, userID string)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID, err := GetUserFromContext(ctx)

		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not authorized to view this page"))
			return
		}

		handler(ctx, userID)
	}
}
