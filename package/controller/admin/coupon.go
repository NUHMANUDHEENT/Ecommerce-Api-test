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
		c.JSON(406, gin.H{
			"status": "Fail",
			"error": "Failed to bind data",
			"code":  406,
		})
		return
	}
	if err := initializer.DB.Create(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"error": "Coupon already exist",
			"code":    500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"message": "New coupon created",
		"data":    couponView,
	})
}

func CouponView(c *gin.Context) {
	var couponView []models.Coupon
	if err := initializer.DB.Find(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"error": "Failed to find coupon details",
			"code":    500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"coupons": couponView,
	})
	couponView = []models.Coupon{}
}

func CouponDelete(c *gin.Context) {
	var couponDelete models.Coupon
	id := c.Param("ID")
	if err := initializer.DB.Where("id=?", id).Delete(&couponDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"error": "Failed to delete coupon",
			"code":    501,
		})
	} else {
		c.JSON(200, gin.H{
			"status": "Success",
			"message": "Coupon deleted succesfully",
		})
	}
}
