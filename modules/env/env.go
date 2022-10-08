package env

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var Env Environ

type Environ struct {
	PORT string
}

func InitEnv() error {
	validator := validator.New()

	godotenv.Load()

	input := Environ{
		PORT: os.Getenv("PORT"),
	}

	if input.PORT == "" {
		input.PORT = "8080"
	}

	err := validator.Struct(input)

	if err != nil {
		return err
	}

	Env = input
	return nil
}
