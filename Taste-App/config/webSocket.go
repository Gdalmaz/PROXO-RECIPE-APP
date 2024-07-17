package config

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// WebSocketHandler handles websocket connections
func WebSocketHandler(c *websocket.Conn) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("Error closing websocket:", err)
		}
	}()

	for {
		// Veri okuma
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// Veri yazma
		if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}
