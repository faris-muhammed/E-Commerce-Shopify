package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

// ================= List Orders =================

func ListOrders(c *gin.Context) {
	requestID := c.GetUint("userid")

	// Fetch order items for the given seller ID and preload the Product
	var orderItems []model.OrderItems
	if err := initializer.DB.Preload("Product").Preload("Order").Where("seller_id = ?", requestID).Find(&orderItems).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch order items",
			"code":   500,
		})
		return
	}

	// Create a map to associate order items with their respective orders
	orderMap := make(map[uint]gin.H)

	// Add order items to their respective orders
	for _, v := range orderItems {
		if _, exists := orderMap[v.OrderId]; !exists {
			// Initialize the order entry in the map if it doesn't exist
			orderMap[v.OrderId] = gin.H{
				"orderId": v.OrderId,
				"userId":  v.UserID,
				"items":   []gin.H{},
			}
		}

		// Append the order item to the order's item list
		orderMap[v.OrderId]["items"] = append(orderMap[v.OrderId]["items"].([]gin.H), gin.H{
			"id":            v.Id,
			"productId":     v.ProductId,
			"productName":   v.Product.ProductName,
			"quantity":      v.Quantity,
			"subTotal":      v.SubTotal,
			"orderStatus":   v.OrderStatus,
			"paymentStatus": v.PaymentStatus,
			"paymentMethod": v.Order.OrderPaymentMethod,
		})
	}

	// Prepare the final response array
	var showOrders []gin.H
	for _, order := range orderMap {
		showOrders = append(showOrders, order)
	}

	// Send the response
	c.JSON(200, gin.H{
		"status": "Success",
		"code":   200,
		"orders": showOrders,
	})
}

//=============== Deliver Order ===============

func DeliverOrder(c *gin.Context) {
	requestID := c.GetUint("userid")
	orderID := c.Param("id")

	// Log the incoming parameters for debugging
	log.Printf("RequestID: %d, OrderID: %s", requestID, orderID)

	var orderItem model.OrderItems
	if err := initializer.DB.Where("id = ?", orderID).First(&orderItem).Error; err != nil {
		// Return an error if the order item is not found
		log.Printf("Error finding order item: %v", err)
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Order item not found",
			"code":   500,
		})
		return
	}
	if orderItem.OrderStatus == "Cancelled" {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Order item already Cancelled",
			"code":   400,
		})
		return
	}

	// Check if the order item is already canceled
	if orderItem.OrderStatus == "Delivered" {
		c.JSON(http.StatusOK, gin.H{
			"status": "Fail",
			"error":  "Order item is already Delivered",
			"code":   400,
		})
		return
	}

	// This is a placeholder check and should be adjusted based on your actual user and order logic
	if orderItem.SellerId != requestID {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "You do not have permission to change the status into Delivered",
			"code":   409,
		})
		return
	}

	// Update the order item's status to "Delivered"
	if err := initializer.DB.Model(&orderItem).Update("OrderStatus", "Delivered").Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to change the status into Delivered",
			"code":   400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Order item Delivered successfully",
		"code":    200,
	})
}

func CancelOrder(c *gin.Context) {
	requestID := c.GetUint("userid")
	orderID := c.Param("id")

	// Log the incoming parameters for debugging
	log.Printf("RequestID: %d, OrderID: %s", requestID, orderID)

	var orderItem model.OrderItems
	if err := initializer.DB.Where("id = ?", orderID).First(&orderItem).Error; err != nil {
		// Return an error if the order item is not found
		log.Printf("Error finding order item: %v", err)
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Order item not found",
			"code":   500,
		})
		return
	}

	if orderItem.OrderStatus == "Delivered" {
		c.JSON(409, gin.H{
			"status": "Success",
			"error":  "Order item already Delivered",
			"code":   409,
		})
		return
	}

	// Check if the order item is already canceled
	if orderItem.OrderStatus == "Cancelled" {
		c.JSON(409, gin.H{
			"status": "success",
			"error":  "Order item already Cancelled",
			"code":   409,
		})
		return
	}

	// This is a placeholder check and should be adjusted based on your actual user and order logic
	if orderItem.SellerId != requestID {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "You do not have permission to change the status into Cancelled",
			"code":   401,
		})
		return
	}

	// Update the order item's status to "Delivered"
	if err := initializer.DB.Model(&orderItem).Update("OrderStatus", "Cancelled").Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to change the status into Cancelled",
			"code":   400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Order item Cancelled successfully",
		"code":    200,
	})
}
