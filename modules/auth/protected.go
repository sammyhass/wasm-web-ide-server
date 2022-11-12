package auth

import "github.com/gin-gonic/gin"

func Protected(handler func(ctx *gin.Context, userID string)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID, err := GetUserFromContext(ctx)

		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		handler(ctx, userID)
	}
}
