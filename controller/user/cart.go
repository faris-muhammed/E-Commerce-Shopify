package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/helper"
	"main.go/initializer"
	"main.go/model"
)

var MaxQuantity = 6

//=========== View Products ===========

func ListProducts(c *gin.Context) {
	var products []model.ProductDetails
	if err := initializer.DB.Where("is_deleted=?", false).Find(&products).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch products",
			"code":   500,
		})
		return
	}

	var productsShow []gin.H
	for _, product := range products {
		productsShow = append(productsShow, gin.H{
			"id":                 product.Id,
			"name":               product.ProductName,
			"price":              product.Price,
			"available_quantity": product.Quantity,
			"size":               product.Size,
		})
	}

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   productsShow,
		"code":   200,
	})
}

// ================ Add to cart ================
func AddCart(c *gin.Context) {
	// Binding the JSON data from the requested URL
	requestID := c.GetUint("userid")
	var cart model.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to parse request body",
			"code":   400,
		})
		return
	}

	// Fetch product from the database
	var product model.ProductDetails
	if err := initializer.DB.First(&product, cart.ProductId).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch product",
			"code":   500,
		})
		return
	}

	// Check if the quantity is less than one
	if cart.Quantity < 1 {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Add quantity greater than zero",
			"code":   400,
		})
		return
	}

	// Check if the product is out of stock
	if product.Quantity < cart.Quantity {
		c.JSON(400, gin.H{
			"error": "Product is out of stock",
			"code":  400,
		})
		return
	}

	// Check if the user already has the product in the cart
	var existingCart model.Cart
	err := initializer.DB.Where("user_id = ? AND product_id = ?", requestID, cart.ProductId).First(&existingCart).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to check existing cart item",
			"code":   500,
		})
		return
	}

	// If the cart item already exists, respond with an appropriate error
	if existingCart.Id != 0 {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Product already in cart",
			"code":   400,
		})
		return
	}

	// If the cart item does not exist, create a new one
	newItem := model.Cart{
		UserId:     requestID,
		ProductId:  cart.ProductId,
		CategoryId: product.CategoryId,
		Quantity:   cart.Quantity,
		Price:      float64(cart.Quantity) * float64(product.Price),
	}

	if err := initializer.DB.Create(&newItem).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to add the product to user's cart",
			"code":   500,
		})
		return
	}

	totalAmount := helper.CalculateTotalAmount(requestID)
	fmt.Println(totalAmount)
	c.JSON(201, gin.H{
		"status":       "Success",
		"message":      "Product added to cart successfully",
		"code":         201,
		"product":      product.ProductName,
		"total_amount": newItem.Price,
	})
}

//========== Edit Cart ==========

func EditCart(c *gin.Context) {
	// Get user ID and product ID from the URL or request parameters
	userID := c.GetUint("userid")
	productID := c.Param("productId")

	// Find the specific cart item based on user ID and product ID
	var existingCart model.Cart
	if err := initializer.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingCart).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Cart item not found",
			"code":   400,
		})
		return
	}

	// Bind the updated cart details from the request body
	var updateData struct {
		Quantity uint `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to bind the data",
			"code":   400,
		})
		return
	}

	// Check if the quantity is zero or negative
	if updateData.Quantity <= 0 {
		// If quantity is zero or negative, delete the cart item
		if err := initializer.DB.Delete(&existingCart).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to delete cart item",
				"code":   500,
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "success",
			"error":  "Cart item deleted successfully",
			"code":   200,
		})
		return
	}
	if updateData.Quantity >= 7 {
		c.JSON(409, gin.H{
			"error": "Add products less than or equal to six",
			"code":  409,
		})
		return
	}

	// Update the cart item quantity
	existingCart.Quantity = uint(updateData.Quantity)

	// Save the updated cart details to the database
	if err := initializer.DB.Save(&existingCart).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to update cart item",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"error":    "Cart item updated successfully",
		"code":     200,
		"id":       existingCart.Id,
		"quantity": existingCart.Quantity,
		"price":    existingCart.Price,
	})
}

// =========== Remove cart ===========

func RemoveCart(c *gin.Context) {
	// Get cart ID from request body
	requestID := c.GetUint("userid") //Bind the data from requested body
	var cart []model.Cart
	if err := initializer.DB.Find(&cart, requestID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch the cart details",
			"code":   400,
		})
		return
	}
	// Delete the cart
	if err := initializer.DB.Where("user_id=?", requestID).Delete(&cart).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete the cart",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"error":  "Cart removed successfully",
		"code":   200,
		"id":     cart,
	})
}

//=============== List Cart =============

func ListCart(c *gin.Context) {
	var carts []model.Cart
	requestID := c.GetUint("userid")
	var cartShow []gin.H
	// Fetch the cart items for the specified user and preload the Product data
	if err := initializer.DB.Preload("Product").Where("user_id = ?", requestID).Find(&carts).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
			"error":  "Failed to fetch cart details",
			"code":   500,
		})
		return
	}

	if len(carts) == 0 {
		c.JSON(409, gin.H{
			"status": "fail",
			"error":  "Please add some items to your cart first.",
			"code":   409,
		})
		return
	}

	var totalAmount float64
	for _, val := range carts {
		amount := val.Product.Price * float64(val.Quantity)
		totalAmount += amount
		cartShow = append(cartShow, gin.H{
			"id":        val.Product.Id,
			"product":   val.Product.ProductName,
			"quantity":  val.Quantity,
			"price":     val.Product.Price,
			"sub_total": amount,
		})
	}

	c.JSON(200, gin.H{
		"status":      "success",
		"code":        200,
		"data":        cartShow,
		"totalAmount": totalAmount,
	})
}
