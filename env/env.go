package env

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var env map[EnvKey]string

var errMissingKey = errors.New("missing required env key: %s")
var errInvalidKey = errors.New("invalid env key: %s")

/*
Get a key from the environment
*/
func Get(key EnvKey) string {
	if val, ok := env[key]; ok {
		return val
	}

	return ""
}

func GetOr(key EnvKey, fallback string) string {
	if val, ok := env[key]; ok {
		return val
	}

	return fallback
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

func Set(key EnvKey, value string) {
	env[key] = value
}

type EnvKey int

const (
	env_none EnvKey = iota
	// START OF ENV KEYS
	// -----------------

	PORT

	// DB
	POSTGRES_HOST
	POSTGRES_PORT
	POSTGRES_USER
	POSTGRES_PASSWORD
	POSTGRES_DB

	// S3
	S3_ACCESS_KEY_ID
	S3_SECRET_ACCESS_KEY
	S3_BUCKET

	JWT_SECRET

	CORS_ALLOW_ORIGIN

	// --------------------
	// END OF ENV KEYS
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
	case S3_ACCESS_KEY_ID:
		return "S3_ACCESS_KEY_ID"
	case S3_SECRET_ACCESS_KEY:
		return "S3_SECRET_ACCESS_KEY"
	case S3_BUCKET:
		return "S3_BUCKET"
	case CORS_ALLOW_ORIGIN:
		return "CORS_ALLOW_ORIGIN"
	default:
		return "INVALID_KEY"
	}
}
