package routers

import (
	"auth/controllers"
	"auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	user := v1.Group("/user")

	user.Post("/signup", controllers.SignUp)
	user.Post("/login", controllers.LogIn)
	user.Put("/", middleware.TokenControl(), controllers.UpdatePassword)
	user.Get("/", middleware.TokenControl(), controllers.LogOut)
	user.Delete("/", middleware.TokenControl(), controllers.DeleteAccount)

}
