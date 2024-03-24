package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func OfferList(c *gin.Context) {
	var offerList []models.Offer
	if err := initializer.DB.Find(&offerList).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to find offers",
		})
	}
	for _, val := range offerList {
		c.JSON(500, gin.H{
			"product id": val.ProductId,
			"offer":      val.SpecialOffer,
			"discount":   val.Discount,
			"valid from": val.ValidFrom,
			"valid to":   val.ValidTo,
		})
	}
}
func OfferAdd(c *gin.Context) {
	var addOffer models.Offer
	err := c.ShouldBindJSON(&addOffer)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to bind data",
		})
		return
	}
	if err := initializer.DB.Create(&addOffer).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to create offer",
		})
		return
	}
	c.JSON(500, gin.H{
		"Message": "New offer created",
	})
}
func OfferDelete(c *gin.Context) {
	var deleteOffer models.Offer
	offerId := c.Param("ID")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to delete offer",
		})
		return
	}
	c.JSON(500, gin.H{
		"Message": "Offer was deleted",
	})
}
