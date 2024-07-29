package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

// ============= List Product =============
func ListProduct(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var product []model.ProductDetails

	if err := initializer.DB.Where("is_deleted=?", false).Find(&product, "seller_id=?", sellerID).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch products",
			"code":   500,
		})
	}
	var showData []gin.H
	for _, v := range product {
		showData = append(showData, gin.H{
			"id":                 v.Id,
			"productName":        v.ProductName,
			"price":              v.Price,
			"Available_quantity": v.Quantity,
			"size":               v.Size,
			"brand":              v.Brand,
			"category":           v.CategoryId,
		})
	}
	c.JSON(200, gin.H{
		"status": "success",
		"code":   200,
		"data":   showData,
	})
}

// ============= Add product ============

func AddProduct(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var product model.ProductDetails
	// Binding the JSON data from the requested URL
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to parse request body",
			"code":   400,
		})
		return
	}
	// Validate product details
	if product.ProductName == "" || product.Price <= 0 {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Product name and price are required",
			"code":   400,
		})
		return
	}

	// Create new Product with object
	newProduct := model.ProductDetails{
		ProductName: product.ProductName,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Size:        product.Size,
		Brand:       product.Brand,
		Barcode:     product.Barcode,
		CategoryId:  product.CategoryId,
		SellerId:    sellerID,
	}

	// Save product details to the database
	if err := initializer.DB.Create(&newProduct).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to add product",
			"code":   500,
		})
		return
	}

	// Response with only the required fields
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "Product added successfully",
		"code":    201,
	})
}

//=========== Edit Product ===========

func EditProduct(c *gin.Context) {
	// Get product ID from the URL
	productID := c.Param("id")
	// Find existing product from the database
	var existingProduct model.ProductDetails
	if err := initializer.DB.First(&existingProduct, productID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Product not found",
			"code":   400,
		})
		return
	}

	// Bind the updated product data from the request body
	if err := c.ShouldBindJSON(&existingProduct); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to bind the data",
			"code":   400,
		})
		return
	}

	// Save the updated product details to the database
	if err := initializer.DB.Save(&existingProduct).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to update product",
			"code":   400,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":      "Success",
		"message":     "Product updated successfully",
		"code":        200,
		"id":          existingProduct.Id,
		"productname": existingProduct.ProductName,
		"price":       existingProduct.Price,
		"quantity":    existingProduct.Quantity,
		"size":        existingProduct.Size,
		"brand":       existingProduct.Brand,
		"barcode":     existingProduct.Barcode,
	})
}

//========== Delete Products ==========

func SoftDeleteProduct(c *gin.Context) {
	var product model.ProductDetails
	requestID := c.Param("id")

	// Find the product by ID
	if err := initializer.DB.First(&product, requestID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Product not found",
			"code":   400,
		})
		return
	}

	// Soft delete the product
	if err := initializer.DB.Model(&product).Update("is_deleted", true).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to soft delete product",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message":     "Product soft deleted successfully",
		"code":        200,
		"id":          product.Id,
		"productname": product.ProductName,
		"price":       product.Price,
		"quantity":    product.Quantity,
		"size":        product.Size,
		"brand":       product.Brand,
		"barcode":     product.Barcode,
	})
}

//========== Recover the Product ==========

func RecoverDeleteProduct(c *gin.Context) {
	var product model.ProductDetails
	requestID := c.Param("id")

	// Find the product by ID
	if err := initializer.DB.First(&product, requestID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Product not found",
			"code":   400,
		})
		return
	}

	// Recover delete the product
	if err := initializer.DB.Model(&product).Update("is_deleted", false).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to recover the product",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message":     "Product recovered successfully",
		"code":        200,
		"id":          product.Id,
		"productname": product.ProductName,
		"price":       product.Price,
		"quantity":    product.Quantity,
		"size":        product.Size,
		"brand":       product.Brand,
		"barcode":     product.Barcode,
	})
}
