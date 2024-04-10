package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// @Summary New coupon create
// @Description  Admin can Create a new coupon with condition and validity
// @Tags /admin/coupon
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.Coupon true "Coupon data"
// @Success 200 {json} JSON "Product created successfully"
// @Failure 406 {json} JSON "Failed to bind data"
// @Failure 500 {json} JSON "Failed to insert Coupon"
// @Router /admin/coupon [POST]
func CouponCreate(c *gin.Context) {
	var couponStore models.Coupon
	err := c.Bind(&couponStore)
	if err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Failed to bind data",
			"code":   406,
		})
		return
	}
	if err := initializer.DB.Create(&couponStore).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Coupon already exist",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "New coupon created",
	})
}
// CouponView godoc
// @Summary Get coupon details
// @Description Get details of all coupons
// @Tags /admin/coupon
// @ID get_all_products
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} JSON  "Products fetched successfully."
// @Failure 500 {json} JSON  "Server error"
// @Router /admin/coupon [get]
func CouponView(c *gin.Context) {
	var couponView []models.Coupon
	var couonShow []gin.H
	if err := initializer.DB.Find(&couponView).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to find coupon details",
			"code":   500,
		})
		return
	}
	for _, v := range couponView {
		couonShow = append(couonShow, gin.H{
			"id":        v.ID,
			"code":      v.Code,
			"condition": v.CouponCondition,
			"validFrom": v.ValidFrom,
			"validTo":   v.ValidTo,
			"discount": v.Discount,
		})
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"coupons": couonShow,
	})
}
// CouponDelete godoc
// @Summary Delete a coupon by ID
// @Description Delete a coupon by its unique identifier
// @Tags /admin/coupon
// @ID deleteCouponByID
// @Param ID path int true "Coupon ID"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} string "Coupon deleted successfully"
// @Failure 500 {json} object "Failed to delete coupon"
// @Router /admin/coupon/{ID} [delete]
func CouponDelete(c *gin.Context) {
	var couponDelete models.Coupon
	id := c.Param("ID")
	if err := initializer.DB.Where("id=?", id).Delete(&couponDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete coupon",
			"code":   501,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "Coupon deleted succesfully",
		})
	}
}
