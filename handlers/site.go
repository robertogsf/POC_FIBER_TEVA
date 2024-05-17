package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/robertogsf/poc_fiber/models"
)

func tienePermisoParaManejarSitios(estado *models.Connections) bool {
	groups_whithout := strings.Trim(estado.AuthGroups, "[]")

	if groups_whithout == "Sites" {
		return true
	}

	authGroups := strings.Split(groups_whithout, ",")

	for _, grupo := range authGroups {

		if strings.TrimSpace(grupo) == "Sites" {
			return true
		}
	}

	return false
}

func PostSite(c *fiber.Ctx) error {
	connectionId := c.Cookies("connectionId")

	estado, encontrado := conexiones[connectionId]

	if !encontrado {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se encontró la conexión",
		})
	}

	if !tienePermisoParaManejarSitios(estado) {
		return models.NewResponse(c, fiber.StatusUnauthorized, "El usuario no tiene permisos para crear un sitio",
			nil, []string{"El usuario no tiene permisos para crear un sitio"})
	}

	sitio := new(models.Site)
	if err := c.BodyParser(sitio); err != nil {
		return models.NewResponse(c, fiber.StatusBadRequest, "No se pudo analizar el cuerpo de la solicitud",
			nil, []string{err.Error()})
	}

	tx := db.Begin()

	if err := tx.Create(&sitio).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al intentar crear el sitio",
			nil, []string{"Hubo un problema al intentar crear el sitio"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Sitio creado con éxito", sitio, nil)
}

func GetSites(c *fiber.Ctx) error {
	connectionId := c.Cookies("connectionId")

	estado, encontrado := conexiones[connectionId]
	if !encontrado {
		return models.NewResponse(c, fiber.StatusUnauthorized, "No se encontró la conexión",
			nil, []string{"No se encontró la conexión"})
	}

	if !tienePermisoParaManejarSitios(estado) {
		return models.NewResponse(c, fiber.StatusUnauthorized, "El usuario no tiene permisos para obtener los sitios",
			nil, []string{"El usuario no tiene permisos para obtener los sitios"})
	}

	var sitios []models.Site
	if err := db.Where("account_id = ?", estado.AccountID).Find(&sitios).Error; err != nil {
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener los sitios",
			nil, []string{"Hubo un problema al obtener los sitios"})
	}

	return models.NewResponse(c, fiber.StatusOK, "Sitios obtenidos con éxito", sitios, nil)
}

func DeleteSite(c *fiber.Ctx) error {

	connectionId := c.Cookies("connectionId")

	estado, encontrado := conexiones[connectionId]
	if !encontrado {
		return models.NewResponse(c, fiber.StatusUnauthorized, "No se encontró la conexión",
			nil, []string{"No se encontró la conexión"})
	}

	if !tienePermisoParaManejarSitios(estado) {
		return models.NewResponse(c, fiber.StatusUnauthorized, "El usuario no tiene permisos para marcar el sitio para borrado",
			nil, []string{"El usuario no tiene permisos para marcar el sitio para borrado"})
	}

	siteId := c.Params("siteId")

	tx := db.Begin()

	var sitio models.Site
	if err := tx.Where("site_id = ? AND account_id = ?", siteId, estado.AccountID).First(&sitio).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Sitio no encontrado",
				nil, []string{"Sitio no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el sitio",
			nil, []string{"Hubo un problema al obtener el sitio"})
	}

	if err := tx.Model(&sitio).Update("record_status", "Deleted").Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al marcar el sitio para borrado",
			nil, []string{"Hubo un problema al marcar el sitio para borrado"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Sitio marcado para borrado con éxito", sitio, nil)
}

func PutSite(c *fiber.Ctx) error {

	connectionId := c.Cookies("connectionId")

	estado, encontrado := conexiones[connectionId]
	if !encontrado {
		return models.NewResponse(c, fiber.StatusUnauthorized, "No se encontró la conexión",
			nil, []string{"No se encontró la conexión"})
	}

	if !tienePermisoParaManejarSitios(estado) {
		return models.NewResponse(c, fiber.StatusUnauthorized, "El usuario no tiene permisos para modificar el sitio",
			nil, []string{"El usuario no tiene permisos para modificar el sitio"})
	}

	siteId := c.Params("siteId")

	tx := db.Begin()

	var sitio models.Site
	if err := tx.Where("id = ? AND account_id = ?", siteId, estado.AccountID).First(&sitio).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Sitio no encontrado",
				nil, []string{"Sitio no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el sitio",
			nil, []string{"Hubo un problema al obtener el sitio"})
	}

	sitioModificado := new(models.Site)
	if err := c.BodyParser(sitioModificado); err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusBadRequest, "No se pudo analizar el cuerpo de la solicitud",
			nil, []string{"No se pudo analizar el cuerpo de la solicitud"})
	}

	if err := tx.Model(&sitio).Updates(sitioModificado).Error; err != nil {
		tx.Rollback()
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al modificar el sitio",
			nil, []string{"Hubo un problema al modificar el sitio"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Sitio modificado con éxito", sitio, nil)
}

func GetSite(c *fiber.Ctx) error {

	connectionId := c.Cookies("connectionId")

	estado, encontrado := conexiones[connectionId]
	if !encontrado {
		return models.NewResponse(c, fiber.StatusUnauthorized, "No se encontró la conexión",
			nil, []string{"No se encontró la conexión"})
	}

	if !tienePermisoParaManejarSitios(estado) {
		return models.NewResponse(c, fiber.StatusUnauthorized, "El usuario no tiene permisos para obtener el sitio",
			nil, []string{"El usuario no tiene permisos para obtener el sitio"})
	}

	siteId := c.Params("siteId")

	tx := db.Begin()

	var sitio models.Site
	if err := tx.Where("id = ? AND account_id = ?", siteId, estado.AccountID).First(&sitio).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return models.NewResponse(c, fiber.StatusNotFound, "Sitio no encontrado",
				nil, []string{"Sitio no encontrado"})
		}
		return models.NewResponse(c, fiber.StatusInternalServerError, "Hubo un problema al obtener el sitio",
			nil, []string{"Hubo un problema al obtener el sitio"})
	}

	tx.Commit()

	return models.NewResponse(c, fiber.StatusOK, "Sitio obtenido con éxito", sitio, nil)
}
