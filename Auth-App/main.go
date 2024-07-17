package main

import (
	"auth/config"
	"auth/database"
	"auth/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Get("/ws", websocket.New(config.WebSocketHandler))
	routers.UserRouter(app)
	app.Listen(":8001")
}
