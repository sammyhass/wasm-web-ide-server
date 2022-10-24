package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUserToContext(ctx *gin.Context, userID string) {
	ctx.Set("user_id", userID)
}

func GetUserFromContext(ctx *gin.Context) string {
	userID, ok := ctx.Get("user_id")

	if !ok {
		return ""
	}

	return userID.(string)
}

func GetUserFromContextOrAbort(ctx *gin.Context) string {
	userID, ok := ctx.Get("user_id")

	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}

	return userID.(string)
}
