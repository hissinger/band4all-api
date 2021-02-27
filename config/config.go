package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	SecretKey string
)

func init() {
	godotenv.Load()
	SecretKey = os.Getenv("SECRET_KEY")
}
