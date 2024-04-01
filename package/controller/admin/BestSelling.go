package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)


func BestProducts(c *gin.Context) {
	var BestProduct []models.Products
	query := c.Query("query")
	fmt.Println("---", query)
	switch query {
	case "product":
		if err := initializer.DB.Table("order_items oi").Select("p.name, p.price , COUNT(oi.quantity) quantity").
			Joins("JOIN products p ON p.id = oi.product_id").
			Group("p.name, p.price").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestProduct).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": err.Error(),
				"code": 500,
			})
			return
		}
		c.JSON(200, BestProduct)
	case "category":
		var BestCategory []models.Category
		if err := initializer.DB.Table("order_items oi").
			Select("c.category_name, COUNT(oi.quantity) AS quantity").
			Joins("JOIN products p ON oi.product_id = p.id").Joins("JOIN categories c ON  c.id=p.category_id").
			Group("c.category_name").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestCategory).Error; err != nil {
			c.JSON(500,gin.H{
				"status":"Fail",
				"message":err,
				"code":500,
			})
			return
		}
		c.JSON(200, BestCategory)
	}

}
