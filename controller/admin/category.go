package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

// =================== List Category ==================
func ListCategory(c *gin.Context) {
	var category []model.Category
	if err := initializer.DB.Where("is_deleted=?", false).Find(&category).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch products",
			"code":   500,
		})
		return
	}
	var showData []gin.H
	for _, v := range category {
		showData = append(showData, gin.H{
			"id":          v.ID,
			"name":        v.Name,
			"description": v.Description,
		})
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": showData,
	})
}

// ============= Add Category =============

func CreateCategory(c *gin.Context) {
	var category model.Category
	// Bind the JSON data from the request body to the category struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Error Binding data",
			"code":   400,
		})
		return
	}
	// Check if required fields are provided
	if category.Name == "" {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Category name is required",
			"code":   400,
		})
		return
	}

	// Create the new category in the database
	if err := initializer.DB.Create(&category).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to create category",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":      "Success",
		"code":        200,
		"id":          category.ID,
		"category":    category.Name,
		"description": category.Description,
	})
}

//=========== Edit Category ===========

func EditCategory(c *gin.Context) {
	// Get category ID from request URL
	categoryID := c.Param("id")
	// Find the existing category from the database
	var existingCategory model.Category
	if err := initializer.DB.First(&existingCategory, categoryID).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Category not found",
			"code":    400,
		})
		return
	}
	// Binding the data
	if err := c.ShouldBindJSON(&existingCategory); err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "Failed to bind the data",
			"code":    400,
		})
		return
	}
	// Save the updated category details to the database
	if err := initializer.DB.Save(&existingCategory).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"message": "Failed to update Category",
			"code":    500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":      "Success",
		"message":     "Category updated successfully",
		"code":        200,
		"id":          existingCategory.ID,
		"category":    existingCategory.Name,
		"description": existingCategory.Description,
	})
}

//============== Delete Category ==============

func DeleteCategory(c *gin.Context) {
	var blockCategory model.Category
	id := c.Param("id")
	err := initializer.DB.First(&blockCategory, "id=?", id)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find Category",
			"code":   500,
		})
		return
	}
	if blockCategory.IsDeleted {
		blockCategory.IsDeleted = true
		c.JSON(200, gin.H{
			"status":        "Success",
			"message":       "Category deleted",
			"code":          200,
			"category_name": blockCategory.Name,
		})
	} else {
		blockCategory.IsDeleted = false
		c.JSON(200, gin.H{
			"status":        "Success",
			"error":         "Category recovered",
			"code":          200,
			"category_name": blockCategory.Name,
		})
	}
	if err := initializer.DB.Save(&blockCategory).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete/recover Category",
			"code":   500,
		})
	}
}
