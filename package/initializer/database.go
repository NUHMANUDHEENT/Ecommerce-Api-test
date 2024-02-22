package initializer

import (
	"log"
	"os"
	"project1/package/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB gorm.DB

func LoadDatabase() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("...........Failed to connect to database..........")
	}
	DB = *db
	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.Products{})
	DB.AutoMigrate(&models.OtpMail{})
	DB.AutoMigrate(&models.Category{})

}
