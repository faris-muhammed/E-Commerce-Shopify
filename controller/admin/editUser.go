package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

//============= Fetch all users =============

func GetAllUsers(c *gin.Context) {
	// Fetch all user details from the database
	var users []model.UserModel
	if err := initializer.DB.Find(&users).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch user details",
			"code":   400,
		})
		return
	}
	var showData []gin.H
	for _, v := range users {
		showData = append(showData, gin.H{
			"id":     v.Id,
			"name":   v.Name,
			"email":  v.Email,
			"mobile": v.Mobile,
			"gender": v.Gender,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   showData,
		"code":   200,
	})
}

// =============== EditUserDetails ===============

func EditUserDetails(c *gin.Context) {
	// Get user ID from the request URL
	userID := c.Param("id")

	// Find the user by ID
	var existingUser model.UserModel
	if err := initializer.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "User not found",
			"code":    400,
		})
		return
	}
	// Binding the data
	if err := c.ShouldBindJSON(&existingUser); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Error binding the data",
			"code":   400,
		})
		return
	}
	// Save updated user details to the database
	if err := initializer.DB.Save(&existingUser).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Failed to update user details",
			"code":    400,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User details updated successfully",
		"code":    200,
		"name":    existingUser.Name,
		"email":   existingUser.Email,
	})
}

//=============== Block User / Unblock User ===============

func BlockUser(c *gin.Context) {

	var blockUser model.UserModel
	id := c.Param("id")
	if err := initializer.DB.First(&blockUser, "id=?", id).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find user",
			"code":   500,
		})
		return
	}
	if blockUser.IsBlocked {
		blockUser.IsBlocked = true
		c.JSON(200, gin.H{
			"status":     "Success",
			"message":    "User's account blocked",
			"code":       200,
			"user_name":  blockUser.Name,
			"user_email": blockUser.Email,
		})
	} else {
		blockUser.IsBlocked = false
		c.JSON(200, gin.H{
			"status":     "Success",
			"error":      "User's account unblocked",
			"code":       200,
			"user_name":  blockUser.Name,
			"user_email": blockUser.Email,
		})
	}
	if err := initializer.DB.Save(&blockUser).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to block/unblock user",
			"code":   500,
		})
	}
}

// =============== Delete User / Recover User ===============

func DeleteUser(c *gin.Context) {
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
