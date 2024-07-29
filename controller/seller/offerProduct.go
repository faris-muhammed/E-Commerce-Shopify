package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/initializer"
	"main.go/model"
)

func OfferProductList(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var offerList []model.OfferProduct
	if err := initializer.DB.Find(&offerList, "seller_id=?", sellerID).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to find offers",
			"code":   404,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"offers": offerList,
	})

}

func OfferProductAdd(c *gin.Context) {
	sellerID := c.GetUint("userid")
	var addOffer model.OfferProduct
	if err := c.Bind(&addOffer); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bind data",
			"err":    err.Error(),
			"code":   400,
		})
		return
	}
	newOffer := model.OfferProduct{
		ProductId:    addOffer.ProductId,
		SpecialOffer: addOffer.SpecialOffer,
		Discount:     addOffer.Discount,
		ValidFrom:    addOffer.ValidFrom,
		ValidTo:      addOffer.ValidTo,
		SellerId:     sellerID,
	}
	if err := initializer.DB.Create(&newOffer).Error; err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "failed to create offer",
			"code":   406,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "New offer created",
	})
}

func OfferProductDelete(c *gin.Context) {
	var deleteOffer model.OfferProduct
	offerId := c.Param("id")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to delete offer",
			"code":   400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Offer was deleted",
	})
}
