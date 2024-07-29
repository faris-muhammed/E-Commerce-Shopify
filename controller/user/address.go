package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

// =================== List Address ==================
func ListAddress(c *gin.Context) {
	var address []model.UserAddress
	if err := initializer.DB.Find(&address).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch products",
			"code":   500,
		})
		return
	}
	var showData []gin.H
	for _, v := range address {
		showData = append(showData, gin.H{
			"id":       v.Id,
			"fullname": v.FullName,
			"mobile":   v.Mobile,
			"address":  v.Address,
			"street":   v.Street,
			"landmark": v.Landmark,
			"pincode":  v.Pincode,
			"city":     v.City,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"code":   200,
		"data":   showData,
	})
}

//=========== Add Address ===========

func AddAddress(c *gin.Context) {
	requestID := c.GetUint("userid")
	var address model.UserAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to parse request body",
			"code":   400,
		})
		return
	}
	// Validate user details
	if address.FullName == "" || address.Mobile == 0 || address.Address == "" || address.Street == "" || address.Landmark == "" || address.Pincode == 0 || address.City == "" {
		c.JSON(409, gin.H{
			"status": "fail",
			"error":  "Fields are empty. Please fill the required fields",
			"code":   409,
		})
		return
	}

	// Create new address with object
	newAddress := model.UserAddress{
		FullName: address.FullName,
		Mobile:   address.Mobile,
		Address:  address.Address,
		Street:   address.Street,
		Landmark: address.Landmark,
		Pincode:  address.Pincode,
		City:     address.City,
		UserId:   requestID,
	}

	// Save address details to the database
	if err := initializer.DB.Create(&newAddress).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to add address",
			"code":   500,
		})
		return
	}

	c.JSON(201, gin.H{
		"status":   "success",
		"error":    "Address added successfully",
		"FullName": address.FullName,
		"Mobile":   address.Mobile,
		"Address":  address.Address,
		"Street":   address.Street,
		"Landmark": address.Landmark,
		"Pincode":  address.Pincode,
		"City":     address.City,
		"code":     201,
	})
}

func EditAddress(c *gin.Context) {
	// Get address ID from the URL
	addressID := c.Param("id")
	// Find existing address from the database
	var existingAddress model.UserAddress
	if err := initializer.DB.First(&existingAddress, addressID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Address not found",
			"code":   400,
		})
		return
	}

	// Bind the updated address details from the request body
	if err := c.ShouldBindJSON(&existingAddress); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to bind the data",
			"code":   400,
		})
		return
	}

	// Save the updated address details to the database
	if err := initializer.DB.Save(&existingAddress).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to update address",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "Success",
		"message":  "Address updated successfully",
		"code":     200,
		"id":       existingAddress.Id,
		"fullname": existingAddress.FullName,
		"mobile":   existingAddress.Mobile,
		"address":  existingAddress.Address,
		"street":   existingAddress.Street,
		"landmark": existingAddress.Landmark,
		"pincode":  existingAddress.Pincode,
		"city":     existingAddress.City,
	})
}

//========== Permanent Delete Address ==========

func DeleteAddress(c *gin.Context) {
	// Get address ID from the URL
	requestID := c.Param("id")
	// Find the address by ID
	var address model.UserAddress
	if err := initializer.DB.First(&address, requestID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Address not found",
			"code":   400,
		})
		return
	}

	// Delete the address
	if err := initializer.DB.Delete(&address).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete address",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "Success",
		"error":    "Address deleted successfully",
		"code":     200,
		"id":       address.Id,
		"fullname": address.FullName,
		"mobile":   address.Mobile,
		"address":  address.Address,
		"street":   address.Street,
		"landmark": address.Landmark,
		"pincode":  address.Pincode,
		"city":     address.City,
	})
}
