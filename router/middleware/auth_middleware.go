package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sammyhass/web-ide/server/auth"
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

	fmt.Println(claims)

	auth.SetUserToContext(ctx, claims["user_id"].(string))
	ctx.Next()

}
