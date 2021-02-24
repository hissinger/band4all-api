package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	jwt "github.com/form3tech-oss/jwt-go"
	jwtware "github.com/gofiber/jwt/v2"
)

var (
	mqttIP    string
	mqttPort  uint16
	username  string
	password  string
	secretKey string
)

type mqtt struct {
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/sessions", newSession)
	v1.Get("/alive", keepAlive)
}

func newSession(c *fiber.Ctx) error {
	mqtt := mqtt{
		IP:       mqttIP,
		Port:     uint16(mqttPort),
		Username: username,
		Password: password,
	}
	return c.JSON(fiber.Map{"mqtt": mqtt})
}

func keepAlive(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Unauthorize error
	if user != username || pass != password {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": encodedToken})
}

func readEnv() {
	mqttIP = os.Getenv("MQTT_IP")
	port, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))
	mqttPort = uint16(port)
	username = os.Getenv("USER_NAME")
	password = os.Getenv("PASSWORD")
	secretKey = os.Getenv("SECRET_KEY")
}

func main() {
	godotenv.Load()

	readEnv()

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
	}))
	app.Post("/login", login)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secretKey),
	}))

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
