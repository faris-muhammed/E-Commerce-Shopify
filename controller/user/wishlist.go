package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/initializer"
	"main.go/model"
)

func GetWishlistItems(c *gin.Context) {
	userID := c.GetUint("userid")
	var wishlist []model.Wishlist
	if err := initializer.DB.Where("user_id = ?", userID).Preload("Product").Find(&wishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch wishlist",
			"code":   400,
		})
		return
	}

	var showData []gin.H
	for _, v := range wishlist {
		showData = append(showData, gin.H{
			"wishlist_id":  v.ID,
			"product_id":   v.ProductID,
			"product_name": v.Product.ProductName,
			"price":        v.Product.Price,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   showData,
		"code":   200,
	})
}

func AddToWishlist(c *gin.Context) {

	userId := c.GetUint("userid")

	var wishlist model.Wishlist
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to bind the data",
			"code":   400,
		})
		return
	}

	// Fetch product from the database to ensure it exists
	var product model.ProductDetails
	if err := initializer.DB.First(&product, wishlist.ProductID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch product",
			"code":   400,
		})
		return
	}

	// Check if the product is already in the wishlist
	var existingWishlist model.Wishlist
	if err := initializer.DB.Where("user_id = ? AND product_id = ?", userId, wishlist.ProductID).First(&existingWishlist).Error; err == nil {
		// Product is already in the wishlist
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "Product is already in wishlist",
			"code":   409,
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An error other than "record not found" occurred
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to check wishlist",
			"code":   500,
		})
		return
	}

	// Add product to wishlist
	newWishlist := model.Wishlist{
		UserID:    userId,
		ProductID: wishlist.ProductID,
	}

	if err := initializer.DB.Create(&newWishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to add to wishlist",
			"code":   400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Product added to wishlist",
		"code":    200,
	})
}

func RemoveProductFromWishlist(c *gin.Context) {
	userID := c.GetUint("userid")
	wishlistId := c.Param("id")

	var wishlistItem model.Wishlist
	if err := initializer.DB.Where("user_id = ? AND id = ?", userID, wishlistId).First(&wishlistItem).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Wishlist item not found",
			"code":   400,
		})
		return
	}

	if err := initializer.DB.Delete(&wishlistItem).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to remove product from wishlist",
			"code":   400,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Product removed from wishlist successfully",
		"code":    200,
	})
}

func RemoveWishlist(c *gin.Context) {
	userID := c.GetUint("userid")

	if err := initializer.DB.Where("user_id = ?", userID).Delete(&model.Wishlist{}).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to remove wishlist",
			"code":   400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Wishlist removed successfully",
		"code":    200,
	})
}
