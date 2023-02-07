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
