package util

import (
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateToken() (string, error) {
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(SECRET_KEY))

	return encodedToken, err
}

func VerficateToken() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(SECRET_KEY),
	})
}
