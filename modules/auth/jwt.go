package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sammyhass/web-ide/server/modules/env"
)

func generateJWTFromUser(
	id string,
) (string, error) {

	return generateJWTFromUserWithClaims(
		id,
		time.Hour*24,
	)
}

func generateJWTFromUserWithClaims(
	id string,
	exp time.Duration,
) (string, error) {

	var claims = jwt.MapClaims{
		"user_id": id,
		"exp":     exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(
		env.Get(env.JWT_SECRET),
	))
}

func VerifyJWT(
	tokenString string,
) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.Get(env.JWT_SECRET)), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}
