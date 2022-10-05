package env

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var Env Environ

type Environ struct {
	POSTGRES_HOST     string `validate:"required"`
	POSTGRES_USER     string `validate:"required"`
	POSTGRES_PASSWORD string `validate:"required"`
	POSTGRES_DB       string `validate:"required"`
	POSTGRES_PORT     string `validate:"required"`
	PORT              string
}

func InitEnv() error {
	validator := validator.New()

	godotenv.Load()

	input := Environ{
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		PORT:              os.Getenv("PORT"),
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
