package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

func SearchProduct(c *gin.Context) {
	searchQuery := c.Query("search")
	sortBy := strings.ToLower(c.DefaultQuery("sort", "a_to_z"))

	// ======== search based query ============
	query := initializer.DB
	if searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}
	// ======== filter products given query =========
	switch sortBy {
	case "price_low_to_high":
		query = query.Order("price asc")
	case "price_high_to_low":
		query = query.Order("price desc")
	case "a_to_z":
		query = query.Order("product_name asc")
	case "z_to_a":
		query = query.Order("product_name desc")
	default:
		query = query.Order("product_name asc")
	}
	var items []model.ProductDetails
	var itemsShow []gin.H
	query.Joins("Category").Find(&items)
	for _, v := range items {
		itemsShow = append(itemsShow, gin.H{
			"name":  v.ProductName,
			"price": v.Price,
			"size":  v.Size,
		})
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   itemsShow,
	})
}
