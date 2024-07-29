package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

//============================== Check Balance ==============================

func WalletBalance(c *gin.Context) {
	userID := c.GetUint("userid")

	var wallet model.Wallet
	if err := initializer.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Wallet not found",
			"code":    400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":         "success",
		"balance in INR": wallet.Balance,
		"code":           200,
	})
}
