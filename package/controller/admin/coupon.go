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
			"error": "Failed to bind data",
		})
		return
	}
	if err := initializer.DB.Create(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Coupon already exist"})
		return
	}
	c.JSON(200, gin.H{
		"message": "New coupon created",
	})
}

func CouponView(c *gin.Context) {
	var couponView []models.Coupon
	if err := initializer.DB.Find(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to find coupon details",
		})
		return
	}
	c.JSON(200, gin.H{
		"coupons": couponView,
	})
	couponView = []models.Coupon{}
}

func CouponDelete(c *gin.Context) {
	var couponDelete models.Coupon
	id := c.Param("ID")
	if err := initializer.DB.Where("id=?", id).Delete(&couponDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete coupon",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Coupon deleted succesfully",
		})
	}
}
