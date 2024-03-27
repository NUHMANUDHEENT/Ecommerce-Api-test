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
		return
	}
	c.JSON(500, gin.H{
		"offers": offerList,
	})

}
func OfferAdd(c *gin.Context) {
	var addOffer models.Offer
	err := c.ShouldBindJSON(&addOffer)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to bind data",
		})
		return
	}
	if err := initializer.DB.Create(&addOffer).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create offer",
		})
		return
	}
	c.JSON(500, gin.H{
		"message": "New offer created",
	})
}
func OfferDelete(c *gin.Context) {
	var deleteOffer models.Offer
	offerId := c.Param("ID")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to delete offer",
		})
		return
	}
	c.JSON(500, gin.H{
		"message": "Offer was deleted",
	})
}
