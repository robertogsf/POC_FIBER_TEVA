// routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertogsf/poc_fiber/handlers"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/auth/login", handlers.PostLogin)
	app.Get("/auth/login", handlers.GetLogin)
	app.Post("/auth", handlers.PostUser)
	app.Get("/auth/:userId", handlers.GetUser)
	app.Put("/auth/:userId", handlers.PutUser)
	app.Delete("/auth/:userId", handlers.DeleteUser)
	app.Post("/sites", handlers.PostSite)
	app.Get("/sites", handlers.GetSites)
	app.Get("/sites/:siteId", handlers.GetSite)
	app.Put("/sites/:siteId", handlers.PutSite)
	app.Delete("/sites/:siteId", handlers.DeleteSite)
}
