package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("Incoming request: %s %s \n", c.Method(), c.Path())
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		ExposeHeaders:    "*",
		AllowCredentials: true,
	}))

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading .env file")
		panic(err)
	}

	app.Post("/api/v1/user/login", forwardToAuth("user/login"))
	app.Post("/api/v1/user/signup", forwardToAuth("user/signup"))
	app.Listen(":80")

}

func forwardToAuth(path string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetURL := os.Getenv("TASTE-APP-AUTH-HOST")+"/api/v1/" + path
		body := bytes.NewReader(c.Body())
		proxyReq ,err := http.NewRequest(c.Method(),targetURL, body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		proxyReq.Header = c.GetReqHeaders()
		client := &http.Client{}
		proxyResp, err := client.Do(proxyReq)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		defer proxyResp.Body.Close()

		bodyBytes, err := io.ReadAll(proxyResp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		
		c.Status(proxyResp.StatusCode)
		for k, v := range proxyResp.Header {
			c.Set(k, v[0])
		}

		return c.Send(bodyBytes)
	}
}
