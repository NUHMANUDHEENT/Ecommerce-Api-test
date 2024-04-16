package initializer

import (
	"log"
	"os"
	"project1/package/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// =================== connect to database ================
func LoadDatabase() {
	dsnUrl := os.Getenv("DSNURL")
	db, err := gorm.Open(postgres.Open(dsnUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("...........Failed to connect to database..........")
	}
	DB = db
	DB.AutoMigrate(&models.Admins{}, &models.Users{}, &models.Products{}, &models.OtpMail{}, &models.Rating{},
		&models.Review{}, &models.Category{}, &models.Address{}, &models.Cart{}, &models.Coupon{},
		&models.Order{}, &models.OrderItems{}, &models.PaymentDetails{}, &models.Wallet{}, &models.Wishlist{}, &models.Offer{})
}
