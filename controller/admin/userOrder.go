package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

func ListOrdersUser(c *gin.Context) {

	userId := c.Param("id")
	var order []model.Order
	if err := initializer.DB.Where("user_id = ?", userId).Find(&order).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Order not found",
			"code":    400,
		})
		return
	}
	var showOrderItems []gin.H
	for _, v := range order {
		showOrderItems = append(showOrderItems, gin.H{
			"user_id":         v.UserId,
			"order_id":        v.Id,
			"coupon_code":     v.CouponCode,
			"order_amount":    v.OrderAmount,
			"payment_method":  v.OrderPaymentMethod,
			"shipping_charge": v.ShippingCharge,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"code":   200,
		"items":  showOrderItems,
	})

}

func ListOrderItemsUser(c *gin.Context) {

	orderId := c.Param("id")

	// Fetch all order items with their related products
	var orderItems []model.OrderItems
	if err := initializer.DB.Preload("Product").Where("order_id = ?", orderId).Find(&orderItems).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Failed to fetch order items",
			"code":    400,
		})
		return
	}

	// Prepare the response array
	var showOrderItems []gin.H
	for _, v := range orderItems {
		showOrderItems = append(showOrderItems, gin.H{
			"id":            v.Id,
			"user_id":       v.UserID,
			"productName":   v.Product.ProductName,
			"quantity":      v.Quantity,
			"subTotal":      v.SubTotal,
			"orderStatus":   v.OrderStatus,
			"paymentStatus": v.PaymentStatus,
		})

	}
	c.JSON(200, gin.H{
		"status": "Success",
		"code":   200,
		"items":  showOrderItems,
	})
}
