package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// AdminOrdersView returns a JSON response with a list of order items for admin view.
// @Summary View admin orders
// @Description Retrieves a list of order items for admin view.
// @Tags Admin/Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} JSON Response "A successful response."
// @Failure 404 {json} JSON  ErrorResponse "An error occurred while processing your request."
// @Router /admin/orders [get]
func AdminOrdersView(c *gin.Context) {
	var orderitems []models.OrderItems
	var orderShow []gin.H
	if err := initializer.DB.Preload("Order").Find(&orderitems).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find orders",
			"code":   404,
		})
		return
	}
	for _, v := range orderitems {
		orderShow = append(orderShow, gin.H{
			"id":          v.Id,
			"orderId":     v.OrderId,
			"productName": v.ProductId,
			"quantity":    v.Quantity,
			"price":       v.SubTotal,
			"status":      v.OrderStatus,
		})
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"data": orderShow,
	})
}

// Admin can cancel any order using  this endpoint.
// @Summary Cancel an order
// @Description Allows the admin to cancel an existing order.
// @Tags Admin/Orders
// @Accept json
// @Produce json
// @Param id path int true "The ID of the order that you want to cancel"
// @Security ApiKeyAuth
// @Success 200 {json} JSON Response "The order has been successfully canceled"
// @Failure 400 {json} JSON  ErrorResponse "An error occurred while cancel the order."
// @Router /admin/ordercancel [patch]
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

// Admin can change order status using  this endpoint.
// @Summary Status update of an order
// @Description Allows the admin to update status of an existing order.
// @Tags Admin/Orders
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "The ID of the order that you want to upfate status"
// @Param status formData string true " New status for the order"
// @Security ApiKeyAuth
// @Success 200 {json} JSON Response "The order status has been changed successfully "
// @Failure 400 {json} JSON  ErrorResponse "An error occurred while updating status of the order."
// @Router /admin/orderstatus [patch]
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
