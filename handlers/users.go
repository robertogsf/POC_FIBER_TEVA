package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/robertogsf/poc_fiber/database"
	"github.com/robertogsf/poc_fiber/models"
)

// Estructura en memoria
var conexiones = database.Conexiones
var db = database.DB

func PostUser(c *fiber.Ctx) error {
	usuario := new(models.User)
	if err := c.BodyParser(usuario); err != nil {
		return models.NewResponse(c, fiber.StatusBadRequest, "No se pudo analizar el cuerpo de la solicitud",
			nil, []string{err.Error()})
	}

	tx := db.Begin()

	if err := tx.Create(&usuario).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al crear el usuario",
			nil, []string{"Hubo un problema al crear el usuario"})
	}

	conexion := &models.Connections{
		UserID: int64(usuario.ID),
	}

	if err := tx.Create(&conexion).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al crear la conexión para el usuario",
			nil, []string{"Hubo un problema al crear la conexión para el usuario"})
	}

	tx.Commit()
	conexiones[conexion.ConnectionID] = conexion

	return models.NewResponse(c, fiber.StatusOK, "Usuario creado con éxito", usuario, nil)
}

func GetUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var usuario models.User
	if err := db.Where("id = ?", userId).First(&usuario).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Usuario no encontrado",
				nil, []string{"Usuario no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el usuario",
			nil, []string{"Hubo un problema al obtener el usuario"})
	}

	return models.NewResponse(c, fiber.StatusOK, "Usuario obtenido con éxito", usuario, nil)
}

func PutUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	tx := db.Begin()

	var usuario models.User
	if err := tx.Where("id = ?", userId).First(&usuario).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Usuario no encontrado",
				nil, []string{"Usuario no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el usuario",
			nil, []string{"Hubo un problema al obtener el usuario"})
	}

	usuarioModificado := new(models.User)
	if err := c.BodyParser(usuarioModificado); err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusBadRequest, "No se pudo analizar el cuerpo de la solicitud",
			nil, []string{"No se pudo analizar el cuerpo de la solicitud"})
	}

	if err := tx.Model(&usuario).Updates(usuarioModificado).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al modificar el usuario",
			nil, []string{"Hubo un problema al modificar el usuario"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Usuario modificado con éxito", usuario, nil)
}

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	tx := db.Begin()

	var usuario models.User
	if err := tx.Where("id = ?", userId).First(&usuario).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Usuario no encontrado",
				nil, []string{"Usuario no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el usuario",
			nil, []string{"Hubo un problema al obtener el usuario"})
	}

	if err := tx.Model(&usuario).Update("record_status", "Deleted").Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al marcar el usuario para borrado",
			nil, []string{"Hubo un problema al marcar el usuario para borrado"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Usuario marcado para borrado con éxito", nil, nil)
}
