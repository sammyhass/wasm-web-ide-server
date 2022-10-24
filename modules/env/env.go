package env

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var Env Environ

type Environ struct {
	PORT string

	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	JWT_SECRET string
}

func InitEnv() error {
	validator := validator.New()

	godotenv.Load()

	input := Environ{
		PORT: os.Getenv("PORT"),

		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		JWT_SECRET:        os.Getenv("JWT_SECRET"),
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
