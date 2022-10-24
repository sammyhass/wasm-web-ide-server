package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/modules/auth"
)

// AuthMiddleware will add the user_id to the context if the user is authenticated
func AuthMiddleware(
	ctx *gin.Context,
) {
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.Next()
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	claims, err := auth.VerifyJWT(tokenString)

	if err != nil {
		ctx.Next()
		return
	}

	auth.SetUserToContext(ctx, claims["user_id"].(string))
	ctx.Next()

}

func Protected(handler func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, err := auth.GetUserFromContext(ctx)

		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		handler(ctx)

	}
}
