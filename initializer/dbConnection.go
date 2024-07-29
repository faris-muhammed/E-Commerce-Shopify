package initializer

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/model"
)

var DB *gorm.DB

func DBconnect() {
	// connecting database
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	DB = db
	if err := DB.AutoMigrate(&model.UserModel{}, &model.AdminModel{}, &model.SellerModel{}, &model.ProductDetails{}, &model.Category{}, &model.UserAddress{}, &model.Cart{}, &model.Order{}, &model.OrderItems{}, model.OTPDetails{}, model.Coupon{}, model.Wishlist{}, model.Wallet{}, model.PaymentDetails{}, model.OfferProduct{}, model.OfferCategory{}); err != nil {
		fmt.Printf("Error migrating database %v: ", err)
	}
}
