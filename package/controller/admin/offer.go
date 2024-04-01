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
			"status": "Fail",
			"error":  "failed to find offers",
			"code":   404,
		})
		return
	}
	c.JSON(500, gin.H{
		"status": "Success",
		"offers": offerList,
	})

}
func OfferAdd(c *gin.Context) {
	var addOffer models.Offer
	err := c.ShouldBindJSON(&addOffer)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bind data",
			"code":   400,
		})
		return
	}
	if err := initializer.DB.Create(&addOffer).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error": "failed to create offer",
			"code":  500,
		})
		return
	}
	c.JSON(500, gin.H{
		"status": "Success",
		"message": "New offer created",
		"data": addOffer,
	})
}
func OfferDelete(c *gin.Context) {
	var deleteOffer models.Offer
	offerId := c.Param("ID")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error": "failed to delete offer",
			"code":  501,
		})
		return
	}
	c.JSON(500, gin.H{
		"status": "Success",
		"message": "Offer was deleted",
	})
}
