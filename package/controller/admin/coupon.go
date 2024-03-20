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
func CouponView(c *gin.Context) {
	var couponView []models.Coupon
	if err := initializer.DB.Find(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to find coupon details",
		})
	} else {
		for _, val := range couponView {
			c.JSON(200, gin.H{
				"Coupon Id":          val.ID,
				"Coupon code":        val.Code,
				"Coupon Discound":    val.Discount,
				"Coupon condition":   val.CouponCondition,
				"Coupon valied from": val.ValidFrom,
				"Coupon valied To":   val.ValidTo,
			})
		}
	}
	couponView = []models.Coupon{}
}
func CouponDelete(c *gin.Context) {
	var couponDelete models.Coupon
	id := c.Param("ID")
	if err := initializer.DB.Where("id=?", id).Delete(&couponDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to delete coupon",
		})
	} else {
		c.JSON(200, gin.H{
			"Message": "coupon deleted succesfully",
		})
	}
}
