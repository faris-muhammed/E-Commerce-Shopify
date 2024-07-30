package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/middleware"
	"main.go/model"
)

var RoleAdmin = "Admin"

func AdminPage(c *gin.Context) {
	// var OrderDetails []models.OrderItems
	var totalSales []model.Order
	var totalAmount float64
	var totalOrder int
	if err := initializer.DB.Find(&totalSales).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"message": "Failed to fetch data",
			"code":    500,
		})
	}
	for _, v := range totalSales {
		totalAmount += v.OrderAmount
		totalOrder += 1
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Welcome to admin page",
		"data":    gin.H{"sales": totalAmount, "orderCount": totalOrder},
		"code":    200,
	})
}

type AdminDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// =============== LOGIN ===============

func AdminLogin(c *gin.Context) {
	var AdminCheck AdminDetails
	var adminStore model.AdminModel

	//Binding the data
	if err := c.ShouldBindJSON(&AdminCheck); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Binding the data",
			"code":   400,
		})
		return
	}

	if AdminCheck.Email == "" || AdminCheck.Password == "" {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Email and password are required",
			"code":   400,
		})
		return
	}
	//Checking the credentials
	if err := initializer.DB.Where("email=? AND password=?", AdminCheck.Email, AdminCheck.Password).First(&adminStore).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Invalid username or password",
			"code":   401,
		})
		return
	}
	//Generating token
	token, err := middleware.GenerateToken(adminStore.Id, adminStore.Email, RoleAdmin)
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to generate JWT token",
			"code":   500,
		})
		return
	}
	//Setting token inside cookie
	c.SetCookie("jwtTokenAdmin", token, int((time.Hour * 1).Seconds()), "", "buynowbazaar.online", false, false)

	c.JSON(200, gin.H{
		"status": "Success",
		"token":  token,
	})
}

// =============== LOGOUT ===============

func AdminLogout(c *gin.Context) {
	// Clear JWT token cookie
	c.SetCookie("jwtTokenAdmin", "", -1, "", "buynowbazaar.online", false, false)

	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Admin logged out successfully",
		"code":    200,
	})
}
