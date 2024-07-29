package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

//============= Fetch all sellers =============

func GetAllSellers(c *gin.Context) {
	// Fetch all seller details from the database
	var sellers []model.SellerModel
	if err := initializer.DB.Find(&sellers).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to fetch user details",
			"code":  400,
		})
		return
	}
	var showData []gin.H
	for _, v := range sellers {
		showData = append(showData, gin.H{
			"id":          v.Id,
			"companyName": v.CompanyName,
			"email":       v.Email,
			"mobile":      v.Mobile,
			"pincode":     v.Pincode,
			"place":       v.Place,
			"gst":         v.Gst,
		})
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": showData,
	})
}

// =============== Edit Seller Details ===============

func EditSellerDetails(c *gin.Context) {
	// Get seller ID from request URL
	sellerID := c.Param("id")
	// Find the seller by ID
	var existingSeller model.SellerModel
	if err := initializer.DB.Where("id = ?", sellerID).First(&existingSeller).Error; err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Seller not found",
			"code":    400,
		})
		return
	}
	//Binding the data
	if err := c.ShouldBindJSON(&existingSeller); err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Error binding the data",
			"code":    400,
		})
		return
	}
	// Save updated seller details to the database
	if err := initializer.DB.Save(&existingSeller).Error; err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Failed to update seller details",
			"code":    400,
		})
		return
	}

	c.JSON(200, gin.H{
		"message":     "Seller details updated successfully",
		"code":        200,
		"companyname": existingSeller.CompanyName,
		"email":       existingSeller.Email,
		"mobile":      existingSeller.Mobile,
		"pincode":     existingSeller.Pincode,
		"place":       existingSeller.Place,
		"gst":         existingSeller.Gst,
	})
}

//=============== Block Seller / Unblock Seller ===============

func BlockSeller(c *gin.Context) {
	var blockSeller model.SellerModel
	id := c.Param("id")
	err := initializer.DB.First(&blockSeller, "id=?", id)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find seller",
			"code":   500,
		})
		return
	}
	if blockSeller.IsBlocked {
		blockSeller.IsBlocked = true
		c.JSON(200, gin.H{
			"status":      "Success",
			"message":     "Seller's account blocked",
			"code":        200,
			"seller_name": blockSeller.CompanyName,
		})
	} else {
		blockSeller.IsBlocked = false
		c.JSON(200, gin.H{
			"status":      "Success",
			"error":       "Seller's account unblocked",
			"code":        200,
			"seller_name": blockSeller.CompanyName,
		})
	}
	if err := initializer.DB.Save(&blockSeller).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to block/unblock seller",
			"code":   500,
		})
	}
}

// // =============== Delete Seller / Recover Seller ===============

func DeleteSeller(c *gin.Context) {
	var deleteSeller model.SellerModel
	id := c.Param("id")
	if err := initializer.DB.First(&deleteSeller, "id=?", id).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Can't find seller",
			"code":   400,
		})
		return
	}
	if err := initializer.DB.Delete(&deleteSeller).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete/recover seller",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Seller deleted successfully",
		"code":    200,
	})
}
