package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertogsf/poc_fiber/database"
	"github.com/robertogsf/poc_fiber/models"
)

func ConnectionMiddleware(c *fiber.Ctx) error {

	connectionId := c.Cookies("connectionId")

	var conexion *models.Connections
	var ok bool
	if conexion, ok = database.Conexiones[connectionId]; !ok {

		var conexionDB models.Connections
		if err := database.DB.Where("connection_id = ?", connectionId).First(&conexionDB).Error; err != nil {

			return models.NewResponse(c, fiber.StatusUnauthorized, "Conexión no encontrada",
				nil, []string{"Conexión no encontrada"})
		}

		conexion = &conexionDB
		database.Conexiones[connectionId] = conexion
	}

	c.Locals("estado", conexion)

	return c.Next()
}
