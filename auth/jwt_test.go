package auth

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {

	token, err := generateJWTFromUser("123")

	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}

	t.Log(token)

}

func TestGenerateValidJWT(t *testing.T) {
	token, err := generateJWTFromUser("123")

	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}

	claims, err := VerifyJWT(token)
	if err != nil {
		t.Error(err)
	}

	if claims["user_id"] != "123" {
		t.Error("user_id is not 123")
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		t.Error("token is expired")
	}

}
func TestGenerateExpiredJWT(t *testing.T) {
	token, err := generateJWTFromUserWithClaims("123", -time.Nanosecond)

	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}

	time.Sleep(time.Nanosecond * 2)

	claims, err := VerifyJWT(token)

	if err == nil {
		t.Error("error is nil")
	}

	if claims != nil {
		t.Error("claims is not nil")
	}

	// check if contains error
	if !strings.Contains(err.Error(), "expired") {
		t.Error("error does not contain expired")
	}

}
