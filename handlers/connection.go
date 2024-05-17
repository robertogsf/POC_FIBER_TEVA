package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/robertogsf/poc_fiber/models"
)

func PostLogin(c *fiber.Ctx) error {
	usuario := new(models.User)
	if err := c.BodyParser(usuario); err != nil {
		return models.NewResponse(c, fiber.StatusBadRequest, "No se pudo analizar el cuerpo de la solicitud",
			nil, []string{"No se pudo analizar el cuerpo de la solicitud"})
	}

	tx := db.Begin()

	var usuarioEncontrado models.User
	if err := tx.Where("email = ? AND password = ?", usuario.Email, usuario.Password).First(&usuarioEncontrado).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusUnauthorized, "Credenciales incorrectas",
				nil, []string{"Credenciales incorrectas"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al verificar las credenciales",
			nil, []string{"Hubo un problema al verificar las credenciales"})
	}

	var conexion models.Connections
	conexion.AccountID = usuarioEncontrado.AccountID
	conexion.UserID = usuarioEncontrado.UserID
	conexion.CreateDate = usuarioEncontrado.CreateDate
	conexion.LastUpdated = usuarioEncontrado.LastUpdated
	conexion.RecordStatus = usuarioEncontrado.RecordStatus
	conexion.Active = usuarioEncontrado.Active
	conexion.AuthMenu = usuarioEncontrado.AuthMenu
	conexion.AuthKeys = usuarioEncontrado.AuthKeys
	conexion.AuthGroups = usuarioEncontrado.AuthGroups
	conexion.Sites = usuarioEncontrado.Sites
	conexion.Clients = usuarioEncontrado.Clients
	conexion.Connected = "True"

	if err := tx.Model(&models.Connections{}).Where("user_id = ?", usuarioEncontrado.ID).Updates(conexion).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al actualizar la conexión",
			nil, []string{"Hubo un problema al actualizar la conexión"})
	}

	if err := tx.Where("user_id = ?", usuarioEncontrado.ID).First(&conexion).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al buscar la conexión actualizada",
			nil, []string{"Hubo un problema al buscar la conexión actualizada"})
	}

	tx.Commit()

	conexiones[conexion.ConnectionID] = &conexion

	c.Cookie(&fiber.Cookie{
		Name:  "connectionId",
		Value: conexion.ConnectionID,
	})
	c.Cookie(&fiber.Cookie{
		Name:  "idioma",
		Value: usuarioEncontrado.I18n,
	})

	c.Locals("estado", &conexion)

	return models.NewResponse(c, fiber.StatusOK, "Inicio de sesión exitoso", nil, nil)
}

func GetLogin(c *fiber.Ctx) error {
	connectionId := c.Cookies("connectionId")

	conexion, encontrado := conexiones[connectionId]

	if !encontrado {
		return models.NewResponse(c, fiber.StatusNotFound, "No se encontró la conexión",
			nil, []string{"No se encontró la conexión"})
	}

	return models.NewResponse(c, fiber.StatusOK, "Conexión obtenida con éxito", conexion, nil)
}
