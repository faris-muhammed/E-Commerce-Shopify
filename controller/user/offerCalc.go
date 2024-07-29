package controller

import (
	"time"

	"main.go/initializer"
	"main.go/model"
)

func OfferDiscountProduct(productId int) float64 {
	var OfferDiscount model.OfferProduct
	var discount float64
	if err := initializer.DB.Where("valid_from < ? AND valid_to > ? AND product_id=?", time.Now(), time.Now(), productId).First(&OfferDiscount).Error; err == nil {
		discount = OfferDiscount.Discount
		return discount
	}
	return 0
}

func OfferDiscountCategory(categoryId int) float64 {
	var OfferDiscount model.OfferCategory
	var discount float64
	if err := initializer.DB.Where("valid_from < ? AND valid_to > ? AND category_id=?", time.Now(), time.Now(), categoryId).First(&OfferDiscount).Error; err == nil {
		discount = OfferDiscount.Discount
		return discount
	}
	return 0
}
