package helper

import (
	"main.go/initializer"
	"main.go/model"
)

func CalculateTotalAmount(requestID uint) float64 {
	var cartItems []model.Cart
	var totalAmount float64

	if err := initializer.DB.Where("user_id = ?", requestID).Find(&cartItems).Error; err != nil {
		return 0
	}

	for _, item := range cartItems {
		totalAmount += float64(item.Price) * float64(item.Quantity)
	}

	return totalAmount
}
