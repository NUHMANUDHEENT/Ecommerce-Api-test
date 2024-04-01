package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AdminOrdersView(c *gin.Context) {
	var ordersitems []models.OrderItems
	if err := initializer.DB.Preload("Order").Find(&ordersitems).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find orders",
			"code":   404,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"orders": ordersitems,
	})
}
func AdminCancelOrder(c *gin.Context) {
	id := c.Param("ID")
	var orderItem models.OrderItems
	tx := initializer.DB.Begin()
	if err := tx.First(&orderItem, id).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find order",
			"code":   404,
		})
		tx.Rollback()
		return
	}
	if orderItem.OrderStatus == "cancelled" {
		c.JSON(202, gin.H{
			"status":  "Warning",
			"message": "this order has been cancelled before.",
		})
		return
	}

	var orderAmount models.Order
	if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to find order details",
			"code":   400,
		})
		tx.Rollback()
		return
	}
	var couponRemove models.Coupon
	if orderAmount.CouponCode != "" {
		if err := initializer.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "can't find coupon code",
				"code":   400,
			})
			tx.Rollback()
		}
		orderAmount.CouponCode = ""
	}
	if couponRemove.CouponCondition > int(orderAmount.OrderAmount) {
		orderAmount.OrderAmount += couponRemove.Discount
		orderAmount.OrderAmount -= float64(orderItem.SubTotal)
		orderAmount.CouponCode = ""
	}

	if err := tx.Save(&orderAmount).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to update order details",
			"code":   400,
		})
		tx.Rollback()
		return
	}
	orderItem.OrderStatus = "cancelled"
	orderItem.OrderCancelReason = "admin cancelled"
	if err := tx.Save(&orderItem).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to change status",
			"code":   500,
		})
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "Order Cancelled",
		"data":    orderItem.OrderStatus,
	})
}

func AdminOrderStatus(c *gin.Context) {
	id := c.Param("ID")
	var orderStatus models.OrderItems
	orderStatusChenge := c.Request.FormValue("status")
	if orderStatusChenge == "" {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Enter the Status",
			"code":   406,
		})
		return
	}
	if err := initializer.DB.First(&orderStatus, id).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find order",
			"code":   404,
		})
		return
	}
	orderStatus.OrderStatus = orderStatusChenge
	initializer.DB.Save(&orderStatus)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "order status changed to  " + orderStatusChenge,
		"data":    orderStatus.OrderStatus,
	})

}
