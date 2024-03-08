package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AdminOrdersView(c *gin.Context) {
	var orders []models.Order

	initializer.DB.Joins("Product").Find(&orders)
	for _, order := range orders {
		c.JSON(200, gin.H{
			"ID":      order.ID,
			"Product": order.Product.Name,
			"Amount":  order.OrderAmount,
			"Status":  order.OrderStatus,
		})
	}
}
func AdminCancelOrder(c *gin.Context) {
	id := c.Param("ID")
	var order models.Order
	if err := initializer.DB.Where("id=?", id).First(&order).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "can't find order",
		})
		return
	}
	order.OrderStatus = "cancelled"
	initializer.DB.Save(&order)
	c.JSON(200, "Order Cancelled.")
}

func AdminOrderStatus(c *gin.Context) {
	id := c.Param("ID")
	var orderStatus models.Order
	orderStatusChenge := c.Request.FormValue("status")
	if orderStatusChenge == "" {
		c.JSON(500, gin.H{
			"Error": "Enter the Status",
		})
		return
	}
	if err := initializer.DB.Where("id=?", id).First(&orderStatus).Error; err != nil {
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
