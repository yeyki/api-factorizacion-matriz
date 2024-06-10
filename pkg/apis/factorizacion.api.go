package apis

import (
	"api-factorizacion-matriz/pkg/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hola mundo!")
	})
	app.Post("/api/factorizar-matriz", services.FactorizarMatriz)
}
