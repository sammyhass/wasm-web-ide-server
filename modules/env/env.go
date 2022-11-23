package env

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var env map[EnvKey]string

/*
Get a key from the environment
*/
func Get(key EnvKey) string {
	val, ok := env[key]
	if !ok {
		return ""
	}

	return val

}

func InitEnv() error {

	godotenv.Load()

	input := make(map[EnvKey]string)

	for key := env_none + 1; key < env_none_final; key++ {
		envVar := os.Getenv(key.String())
		isOpt := strings.HasPrefix(key.String(), "OPT_")

		if envVar == "" && !isOpt {
			log.Fatalf("Missing required key %s", key.String())
		} else if envVar != "" {
			input[key] = envVar
		}
	}

	env = input
	return nil
}

type EnvKey int

const (
	env_none EnvKey = iota
	// START OF ENV KEYS

	PORT
	POSTGRES_HOST
	POSTGRES_PORT
	POSTGRES_USER
	POSTGRES_PASSWORD
	POSTGRES_DB

	JWT_SECRET
	COOKIE_SECRET

	S3_ACCESS_KEY_ID
	S3_SECRET_ACCESS_KEY
	S3_BUCKET

	// END OF ENV KEYS - final used for loading in everything
	env_none_final
)

func (e EnvKey) String() string {
	switch e {
	case PORT:
		return "PORT"
	case POSTGRES_HOST:
		return "POSTGRES_HOST"
	case POSTGRES_PORT:
		return "POSTGRES_PORT"
	case POSTGRES_USER:
		return "POSTGRES_USER"
	case POSTGRES_PASSWORD:
		return "POSTGRES_PASSWORD"
	case POSTGRES_DB:
		return "POSTGRES_DB"
	case JWT_SECRET:
		return "JWT_SECRET"
	case COOKIE_SECRET:
		return "COOKIE_SECRET"
	case S3_ACCESS_KEY_ID:
		return "S3_ACCESS_KEY_ID"
	case S3_SECRET_ACCESS_KEY:
		return "S3_SECRET_ACCESS_KEY"
	case S3_BUCKET:
		return "S3_BUCKET"
	default:
		return "INVALID_KEY"
	}
}
