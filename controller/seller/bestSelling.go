package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

func BestSelling(c *gin.Context) {
	var BestProduct []model.ProductDetails
	var BestList []gin.H
	query := c.Query("type")
	switch query {
	case "product":
		if err := initializer.DB.Table("order_items oi").Select("p.product_name, p.price , COUNT(oi.quantity) quantity").
			Joins("JOIN product_details p ON p.id = oi.product_id").
			Group("p.product_name, p.price").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestProduct).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": "Failed to fetch data",
				"code":    500,
			})
			return
		}
		for _, v := range BestProduct {
			BestList = append(BestList, gin.H{
				"productName": v.ProductName,
				"salesVolume": v.Quantity,
			})
		}

	case "category":
		var BestCategory []model.Category
		if err := initializer.DB.Table("order_items oi").
			Select("c.name, COUNT(oi.quantity) AS quantity").
			Joins("JOIN product_details p ON oi.product_id = p.id").Joins("JOIN categories c ON  c.id=p.category_id").
			Group("c.name").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestCategory).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": "Failed to fetch data",
				"code":    500,
			})
			return
		}
		for _, v := range BestCategory {
			BestList = append(BestList, gin.H{
				"categoryName": v.Name,
			})
		}
	}
	c.JSON(200, gin.H{
		"data":   BestList,
		"status": "Success",
		"code":   200,
	})
}
