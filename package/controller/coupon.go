package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func CouponStore(c *gin.Context) {
	var couponView models.Coupon
	err := c.ShouldBindJSON(&couponView)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Failed to bind data",
		})
	} else {
		if err := initializer.DB.Create(&couponView).Error; err != nil {
			c.JSON(500, gin.H{
				"Error": "Coupon already exist"})
		} else {
			c.JSON(200, gin.H{
				"message": "New coupon created",
			})
		}
	}
}
