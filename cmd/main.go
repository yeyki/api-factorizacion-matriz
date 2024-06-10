package main

import (
	"api-factorizacion-matriz/pkg/apis"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inicializar Fiber
	app := fiber.New()

	// Definir Rutas
	apis.SetupRoutes(app)

	// Iniciar el Servidor
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
