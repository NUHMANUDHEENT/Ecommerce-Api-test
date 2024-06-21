package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)
// OfferList godoc
// @Summary Get a list of offers
// @Description Retrieve a list of all available offers
// @Tags Admin/Offer
// @ID getOfferList
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} JSON "OK"
// @Failure 400 {string} string error message
// @Router /admin/offer [get]
func OfferList(c *gin.Context) {
	var offerList []models.Offer
	if err := initializer.DB.Find(&offerList).Error; err != nil {
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
// OfferAdd godoc
// @Summary Add a new offer
// @Description Add a new offer to the system
// @Tags Admin/Offer
// @ID addOffer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param offer body models.Offer true "Offer details"
// @Success 200 {json}  JSON "New Offer Created"
// @Failure 400 {json}  ErrorResponse "Failed to create offer"
// @Router /admin/offer [post]
func OfferAdd(c *gin.Context) {
	var addOffer models.Offer
	err := c.Bind(&addOffer)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bind data",
			"code":   400,
		})
		return
	}
	if err := initializer.DB.Create(&addOffer).Error; err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error": "failed to create offer",
			"code":  406,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"message": "New offer created",
	})
}
// OfferDelete godoc
// @Summary Delete an offer by ID
// @Description Delete an offer from the system by its unique identifier
// @Tags Admin/Offer
// @ID deleteOffer
// @Produce json
// @Security ApiKeyAuth
// @Param ID path int true "Offer ID"
// @Success 200 {json}  string  "Deleted Successfully"
// @Failure 400 {json}     ErrorResponse "Failed to delete offer"
// @Router /admin/offer/{ID} [delete]
func OfferDelete(c *gin.Context) {
	var deleteOffer models.Offer
	offerId := c.Param("ID")
	if err := initializer.DB.Where("id=?", offerId).Delete(&deleteOffer).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error": "failed to delete offer",
			"code":  400,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"message": "Offer was deleted",
	})
}
