package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/sammyhass/web-ide/server/modules/env"
)

func GenerateJWTFromClaims(
	claims map[string]interface{},
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	return token.SignedString([]byte(
		env.Env.JWT_SECRET,
	))
}

func VerifyJWT(
	tokenString string,
) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.Env.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
