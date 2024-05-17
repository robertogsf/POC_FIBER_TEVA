package main

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/robertogsf/poc_fiber/database"
	"github.com/robertogsf/poc_fiber/middlewares"
	"github.com/robertogsf/poc_fiber/routes"
)

func main() {
	app := fiber.New()
	app.Use(middlewares.ConnectionMiddleware)
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
