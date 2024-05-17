package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/robertogsf/poc_fiber/models"
	"github.com/robertogsf/poc_fiber/tools"
)

var (
	DB         *gorm.DB
	Conexiones map[string]*models.Connections = make(map[string]*models.Connections)
)

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	host := os.Getenv("DB_HOST")
	tools.IsEmpty(host, "DB_HOST")

	port := os.Getenv("DB_PORT")
	tools.IsEmpty(port, "DB_PORT")

	user := os.Getenv("DB_USER")
	tools.IsEmpty(user, "DB_USER")

	dbname := os.Getenv("DB_NAME")
	tools.IsEmpty(dbname, "DB_NAME")

	password := os.Getenv("DB_PASSWORD")
	tools.IsEmpty(password, "DB_PASSWORD")
	fmt.Println(host, port, user, dbname, password)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	DB, err = gorm.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Connections{})
	DB.AutoMigrate(&models.Site{})

}
