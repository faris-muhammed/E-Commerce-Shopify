package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

func CouponView(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var couponView []model.Coupon
	if err := initializer.DB.Where("seller_id=?", sellerID).Find(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to find coupon details",
			"code":   500,
		})
		return
	}
	var couponShow []gin.H
	for _, v := range couponView {
		couponShow = append(couponShow, gin.H{
			"id":              v.Id,
			"code":            v.Code,
			"created":         v.CreatedDate,
			"expiration_date": v.ExpirationDate,
			"discount":        v.Discount,
		})
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"coupons": couponShow,
		"code":    200,
	})
}

func CouponCreate(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var couponStore model.Coupon
	if err := c.ShouldBindJSON(&couponStore); err != nil {
		c.JSON(406, gin.H{
			"status":  "Fail",
			"message": "Failed to bind data",
			"code":    406,
		})
		return
	}
	newCoupon := model.Coupon{
		Code:           couponStore.Code,
		Discount:       couponStore.Discount,
		MaxDiscount:    couponStore.MaxDiscount,
		MinPurchase:    couponStore.MinPurchase,
		CreatedDate:    couponStore.CreatedDate,
		ExpirationDate: couponStore.ExpirationDate,
		SellerId:       sellerID,
	}
	if err := initializer.DB.Create(&newCoupon).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Coupon already exist",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "New coupon created",
		"code":    200,
	})
}

func CouponDelete(c *gin.Context) {
	var couponDelete model.Coupon
	id := c.Param("id")
	if err := initializer.DB.Where("id=?", id).Delete(&couponDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete coupon",
			"code":   500,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "Coupon deleted succesfully",
			"code":    200,
		})
	}
}
