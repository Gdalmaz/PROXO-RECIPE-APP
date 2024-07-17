package routers

import (
	"proxo-go-application/controllers"
	"proxo-go-application/middleware"

	"github.com/gofiber/fiber/v2"
)

func FoodRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	food := v1.Group("/food")

	food.Post("/", middleware.TokenControl(), controllers.AddTaste)
	food.Put("/:id", middleware.TokenControl(), controllers.UpdateTaste)
	food.Delete("/:id", middleware.TokenControl(), controllers.DeleteTaste)
	food.Get("/", middleware.TokenControl(), controllers.GetAllTaste)
	food.Get("/your-foods", middleware.TokenControl(), controllers.GetAllYourTaste)
	food.Get("/gettaste", middleware.TokenControl(), controllers.GetClickTaste)
	food.Get("/getpopulartaste", middleware.TokenControl(), controllers.PopularTaste)
	food.Get("/search", middleware.TokenControl(), controllers.SearchHandler)
}
