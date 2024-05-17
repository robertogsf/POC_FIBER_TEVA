package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/robertogsf/poc_fiber/tools"
	"golang.org/x/text/language"
)

type RecordStatus string

const (
	Active RecordStatus = "app.status.Active"
)

type User struct {
	gorm.Model
	AccountID      int64     `gorm:"type:bigint" validate:"required"`
	UserID         int64     `gorm:"type:bigserial" validate:"required"`
	CreateDate     time.Time `validate:"required"`
	LastUpdated    time.Time
	RecordStatus   RecordStatus
	Active         RecordStatus
	Identification string `validate:"omitempty,len=50"`
	Password       string `validate:"omitempty,len=50"`
	CompanyName    string `validate:"omitempty,len=50"`
	FirstName      string `validate:"omitempty,len=50"`
	LastName       string `validate:"omitempty,len=50"`
	Email          string `validate:"omitempty,email,len=50"`
	Phone          string `validate:"omitempty,len=50"`
	EmergencyPhone string `validate:"omitempty,len=50"`
	I18n           string `validate:"omitempty,len=50"`
	Address        string `validate:"omitempty,len=200"`
	AuthMenu       string `gorm:"type:varchar;default:'[]'"`
	AuthKeys       string `gorm:"type:varchar;default:'[]'"`
	AuthGroups     string `gorm:"type:varchar;default:'[]'"`
	Sites          string `gorm:"type:varchar;default:'[]'"`
	Clients        string `gorm:"type:varchar;default:'[]'"`
}

type Connections struct {
	gorm.Model
	ConnectionID string `validate:"required,len=50"`
	AccountID    int64
	UserID       int64
	CreateDate   time.Time `validate:"required"`
	LastUpdated  time.Time
	RecordStatus RecordStatus
	Active       RecordStatus
	Connected    string `validate:"omitempty,len=50"`
	Disconnected string `validate:"omitempty,len=50"`
	UserData     string `gorm:"type:varchar;default:'{}'"`
	AccountData  string `gorm:"type:varchar;default:'{}'"`
	AuthMenu     string `gorm:"type:varchar;default:'[]'"`
	AuthKeys     string `gorm:"type:varchar;default:'[]'"`
	AuthGroups   string `gorm:"type:varchar;default:'[]'"`
	Sites        string `gorm:"type:varchar;default:'[]'"`
	Clients      string `gorm:"type:varchar;default:'[]'"`
}

type Site struct {
	gorm.Model
	AccountID         int64     `gorm:"type:bigint" validate:"required"`
	SiteID            int64     `gorm:"type:bigserial" validate:"required"`
	SiteName          string    `validate:"required,len=50"`
	CreateDate        time.Time `validate:"required"`
	LastUpdated       time.Time
	RecordStatus      RecordStatus
	Active            RecordStatus
	Des               string `validate:"omitempty,len=200"`
	Description       string `validate:"omitempty"`
	OperateBy         string `validate:"omitempty,len=50"`
	Logo              string `validate:"omitempty"`
	RulesDocuments    string `gorm:"type:varchar;default:'{}'"`
	ServicesAmenities string `gorm:"type:varchar;default:'{}'"`
	Type              string `gorm:"default:'Yard'" validate:"omitempty,len=50"`
	Email             string `validate:"omitempty,email,len=50"`
	Phone             string `validate:"omitempty,len=50"`
	Address           string `validate:"omitempty,len=200"`
	Website           string `validate:"omitempty,len=200"`
	Geolocation       string `validate:"omitempty,len=200"`
}

func ValidateModel(v *validator.Validate, model interface{}) error {
	return v.Struct(model)
}

type Response struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Method   string      `json:"method"`
	Resource string      `json:"resource"`
	Errors   []string    `json:"errors"`
	Data     interface{} `json:"data"`
}

func NewResponse(c *fiber.Ctx, code int, message string, data interface{}, errors []string) error {

	var idioma language.Tag

	idiomaCookie := c.Cookies("idioma")
	if idiomaCookie != "" {
		idioma = language.Make(idiomaCookie)
	} else {
		idioma = language.English
	}
	if idioma != language.Spanish {
		message = tools.I18n(idioma, message)

		for i, err := range errors {
			errors[i] = tools.I18n(idioma, err)
		}
	}

	response := &Response{
		Code:     code,
		Message:  message,
		Method:   c.Method(),
		Resource: c.Path(),
		Errors:   errors,
		Data:     data,
	}

	return c.JSON(response)
}
