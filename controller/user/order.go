package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/initializer"
	"main.go/model"
)

//=================== List orders ======================

func ListOrders(c *gin.Context) {
	requestID := c.GetUint("userid")
	var order []model.Order
	if err := initializer.DB.Where("user_id = ?", requestID).Find(&order).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"error":   "Order not found or you do not have permission to view it",
			"details": err.Error(),
			"code":    http.StatusForbidden,
		})
		return
	}
	var showOrderItems []gin.H
	for _, v := range order {
		showOrderItems = append(showOrderItems, gin.H{
			"order_id":       v.Id,
			"coupon_code":    v.CouponCode,
			"order_amount":   v.OrderAmount,
			"payment_method": v.OrderPaymentMethod,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   http.StatusOK,
		"items":  showOrderItems,
	})

}

func ListOrderItems(c *gin.Context) {
	requestID := c.GetUint("userid")
	orderId := c.Param("id")

	// Verify the order belongs to the user
	var order model.Order
	if err := initializer.DB.Where("id = ? AND user_id = ?", orderId, requestID).First(&order).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"error":   "Order not found or you do not have permission to view it",
			"details": err.Error(),
			"code":    http.StatusForbidden,
		})
		return
	}

	// Fetch all order items with their related products
	var orderItems []model.OrderItems
	if err := initializer.DB.Preload("Product").Where("order_id = ?", orderId).Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"error":   "Failed to fetch order items",
			"details": err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// Prepare the response array
	var showOrderItems []gin.H
	for _, v := range orderItems {
		showOrderItems = append(showOrderItems, gin.H{
			"id":            v.Id,
			"productName":   v.Product.ProductName,
			"quantity":      v.Quantity,
			"subTotal":      v.SubTotal,
			"orderStatus":   v.OrderStatus,
			"paymentStatus": v.PaymentStatus,
		})

	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   http.StatusOK,
		"items":  showOrderItems,
	})
}

//==================== Cancel Order ====================

func CancelOrderItem(c *gin.Context) {
	requestID := c.GetUint("userid")
	orderID := c.Param("id")

	// Log the incoming parameters for debugging
	log.Printf("RequestID: %d, OrderID: %s", requestID, orderID)

	var orderItem model.OrderItems
	if err := initializer.DB.Where("id = ?", orderID).Preload("Order").First(&orderItem).Error; err != nil {
		// Return an error if the order item is not found
		log.Printf("Error finding order item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Order item not found",
			"err":   err.Error(),
			"code":  http.StatusInternalServerError,
		})
		return
	}
	if orderItem.OrderStatus == "Delivered" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "success",
			"message": "Order item already Delivered",
		})
		return
	}

	// Check if the order item is already canceled
	if orderItem.OrderStatus == "Cancelled" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Order item is already cancelled",
		})
		return
	}

	// Optional: Validate if the user has permission to cancel this order item
	// This is a placeholder check and should be adjusted based on your actual user and order logic
	if orderItem.UserID != requestID {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "Fail",
			"error":  "You do not have permission to cancel this order item",
			"code":   http.StatusForbidden,
		})
		return
	}

	// Update the order item's status to "Cancelled"
	if err := initializer.DB.Model(&orderItem).Update("OrderStatus", "Cancelled").Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Fail",
			"error":  "Failed to cancel order item",
			"code":   http.StatusUnauthorized,
		})
		return
	}

	// Add the quantity back to the product
	var product model.ProductDetails
	if err := initializer.DB.Where("id = ?", orderItem.ProductId).First(&product).Error; err != nil {
		log.Printf("Error finding product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Product not found",
			"err":   err.Error(),
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// Update the product quantity
	product.Quantity += orderItem.Quantity

	if err := initializer.DB.Save(&product).Error; err != nil {
		log.Printf("Error updating product quantity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product quantity",
			"err":   err.Error(),
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// Update the order amount
	orderItem.Order.OrderAmount -= product.Price
	if err := initializer.DB.Save(&orderItem.Order).Error; err != nil {
		log.Printf("Error updating order amount: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order amount",
			"err":   err.Error(),
			"code":  http.StatusInternalServerError,
		})
		return
	}
	if orderItem.PaymentStatus == "success" {
		var wallet model.Wallet
		if err := initializer.DB.Where("user_id = ?", requestID).First(&wallet).Error; err != nil {
			// If wallet not found, create a new one
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = model.Wallet{
					UserID:  requestID,
					Balance: 0, // Initial balance can be zero or any default value
				}
				if err := initializer.DB.Create(&wallet).Error; err != nil {
					log.Printf("Error creating new wallet: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"status": "Fail",
						"error":  "Failed to create wallet",
						"code":   http.StatusInternalServerError,
					})
					return
				}
			} else {
				log.Printf("Error finding wallet: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Wallet not found",
					"err":   err.Error(),
					"code":  http.StatusInternalServerError,
				})
				return
			}
		}

		wallet.Balance += product.Price

		if err := initializer.DB.Save(&wallet).Error; err != nil {
			log.Printf("Error updating wallet balance: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Fail",
				"error":  "Failed to update wallet balance",
				"code":   http.StatusInternalServerError,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Order item canceled successfully",
		"code":    http.StatusOK,
	})
}
