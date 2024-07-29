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
			"status":  "Fail",
			"message": "Failed to fetch products",
			"code":    500,
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

func OfferCategoryList(c *gin.Context) {
	sellerID := c.GetUint("userid")

	var offerList []model.OfferCategory
	if err := initializer.DB.Find(&offerList, "seller_id=?", sellerID).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "failed to find offers",
			"code":    400,
		})
		return
	}
	var showData []gin.H
	for _, v := range offerList {
		showData = append(showData, gin.H{
			"id":         v.Id,
			"offer":      v.SpecialOffer,
			"discount":   v.Discount,
			"valid_from": v.ValidFrom,
			"valid_to":   v.ValidTo,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"offers": showData,
		"code":   200,
	})

}

func OfferCategoryAdd(c *gin.Context) {
	sellerID := c.GetUint("userid")

	var addOffer model.OfferCategory
	if err := c.Bind(&addOffer); err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"message": "Failed to bind data",
			"code":    400,
		})
		return
	}
	newOffer := model.OfferCategory{
		CategoryId:   addOffer.CategoryId,
		SpecialOffer: addOffer.SpecialOffer,
		Discount:     addOffer.Discount,
		ValidFrom:    addOffer.ValidFrom,
		ValidTo:      addOffer.ValidTo,
		SellerId:     sellerID,
	}
	if err := initializer.DB.Create(&newOffer).Error; err != nil {
		c.JSON(406, gin.H{
			"status":  "Fail",
			"message": "Failed to create offer",
			"code":    406,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "New offer created",
		"code":    200,
	})
}

func OfferCategoryDelete(c *gin.Context) {
	var deleteOffer model.OfferCategory
	offerId := c.Param("id")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "failed to delete offer",
			"code":    400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Offer was deleted",
		"code":    200,
	})
}
