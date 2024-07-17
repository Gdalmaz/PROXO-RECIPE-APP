package main

import (
	"proxo-go-application/config"
	"proxo-go-application/database"
	"proxo-go-application/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {

	database.Connect()
	config.ConnectRedis()
	app := fiber.New()
	app.Get("/ws", websocket.New(config.WebSocketHandler))
	routers.FoodRouter(app)
	app.Listen(":8000")
}
