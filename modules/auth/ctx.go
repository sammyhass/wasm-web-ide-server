package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetUserToContext(ctx *gin.Context, userID string) {
	ctx.Set("user_id", userID)
}

func GetUserFromContext(ctx *gin.Context) (string, error) {
	userID, ok := ctx.Get("user_id")

	if !ok {
		return "", fmt.Errorf("not authenticated")
	}

	return userID.(string), nil
}

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
