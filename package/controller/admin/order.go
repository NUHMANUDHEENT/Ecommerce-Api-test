package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AdminOrdersView(c *gin.Context) {
	var ordersitems []models.OrderItems
	if err := initializer.DB.Preload("Order").Find(&ordersitems).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "can't find orders",
		})
		return
	}
	c.JSON(200, gin.H{
		"orders": ordersitems,
	})
}
func AdminCancelOrder(c *gin.Context) {
	id := c.Param("ID")
	var orderItem models.OrderItems
	tx := initializer.DB.Begin()
	if err := tx.First(&orderItem, id).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "can't find order",
		})
		tx.Rollback()
		return
	}
	if orderItem.OrderStatus == "cancelled" {
		c.JSON(202, gin.H{
			"message": "product already cancelled",
		})
		return
	}
	orderItem.OrderStatus = "cancelled"
	orderItem.OrderCancelReason = "admin cancelled"
	if err := tx.Save(&orderItem).Error; err != nil {
		c.JSON(500, "Failed to update status")
		tx.Rollback()
		return
	}

	var orderAmount models.Order
	if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "failed to find order details",
		})
		tx.Rollback()
		return
	}
	var couponRemove models.Coupon
	if orderAmount.CouponCode != "" {
		if err := initializer.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
			c.JSON(400, gin.H{
				"error": "can't find coupon code",
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
			"error": "failed to update order details",
		})
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(201, gin.H{
		"message": "Order Cancelled",
	})
}

func AdminOrderStatus(c *gin.Context) {
	id := c.Param("ID")
	var orderStatus models.OrderItems
	orderStatusChenge := c.Request.FormValue("status")
	if orderStatusChenge == "" {
		c.JSON(500, gin.H{
			"error": "Enter the Status",
		})
		return
	}
	if err := initializer.DB.First(&orderStatus, id).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "can't find order",
		})
		return
	}
	orderStatus.OrderStatus = orderStatusChenge
	initializer.DB.Save(&orderStatus)
	c.JSON(200, gin.H{
		"message": "order status changed to  " + orderStatusChenge,
	})

}
