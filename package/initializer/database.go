package initializer

import (
	"log"
	"os"
	"project1/package/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("...........Failed to connect to database..........")
	}

	DB = db

	// Automatically migrate the schema
	DB.AutoMigrate(&models.Admins{}, &models.Users{}, &models.Products{}, &models.OtpMail{}, &models.Rating{},
		&models.Review{}, &models.Category{}, &models.Address{}, &models.Cart{}, &models.Coupon{},
		&models.Order{}, &models.OrderItems{}, &models.PaymentDetails{}, &models.Wallet{}, &models.Wishlist{}, &models.Offer{})
}
