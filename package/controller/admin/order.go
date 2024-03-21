package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AdminOrdersView(c *gin.Context) {
	var ordersitems []models.OrderItems
	initializer.DB.Preload("Order").Find(&ordersitems)
	for _, orderitem := range ordersitems {
		c.JSON(200, gin.H{
			"Order item id":  orderitem.Id,
			"order id":       orderitem.OrderId,
			"total Amount":   orderitem.SubTotal,
			"user id":        orderitem.Order.UserId,
			"payment method": orderitem.Order.OrderPaymentMethod,
			"order date":     orderitem.Order.OrderDate,
		})
	}
}
func AdminCancelOrder(c *gin.Context) {
	id := c.Param("ID")
	var orderItem models.OrderItems
	tx := initializer.DB.Begin()
	if err := tx.First(&orderItem, id).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "can't find order",
		})
		tx.Rollback()
		return
	}
	if orderItem.OrderStatus == "cancelled" {
		c.JSON(202, gin.H{
			"Message": "product already cancelled",
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
			"Error": "failed to find order details",
		})
		tx.Rollback()
		return
	}
	var couponRemove models.Coupon
	if orderAmount.CouponCode != "" {
		if err := initializer.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
			c.JSON(400, gin.H{
				"Error": "can't find coupon code",
			})
			tx.Rollback()
		}
		orderAmount.CouponCode = ""
	}
	newAmount := 0.0
	newAmount = float64(orderItem.SubTotal) + couponRemove.Discount
	orderAmount.OrderAmount -= newAmount

	if err := tx.Save(&orderAmount).Error; err != nil {
		c.JSON(400, gin.H{
			"Error": "failed to update order details",
		})
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(201, gin.H{
		"Message": "Order Cancelled",
	})
}

func AdminOrderStatus(c *gin.Context) {
	id := c.Param("ID")
	var orderStatus models.OrderItems
	orderStatusChenge := c.Request.FormValue("status")
	if orderStatusChenge == "" {
		c.JSON(500, gin.H{
			"Error": "Enter the Status",
		})
		return
	}
	if err := initializer.DB.First(&orderStatus, id).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "can't find order",
		})
		return
	}
	orderStatus.OrderStatus = orderStatusChenge
	initializer.DB.Save(&orderStatus)
	c.JSON(200, gin.H{
		"Message": "order status changed to  " + orderStatusChenge,
	})

}
