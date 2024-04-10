package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// @Summary Best selling products
// @Description Fetch  the best selling products from database
// @Tags admin
// @Accept json
// @Produce json
// @Secure ApiKeyAuth
// @Param type query string true "Type of search: 'product' or 'category'"
// @Success 200 {json} JSON "User was deleted"
// Failure 404 {json} JSON  "ErrorResponse"
// @Router /admin/bestselling [get]
func BestSelling(c *gin.Context) {
	var BestProduct []models.Products
	var BestList []gin.H
	query := c.Query("type")
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
				"code":    500,
			})
			return
		}
		for _, v := range BestProduct {
			BestList = append(BestList, gin.H{
				"productName":  v.Name,
				"salesVolume":  v.Quantity,
				"averagePrice": float64(v.Price) / float64(v.Quantity),
			})
		}

	case "category":
		var BestCategory []models.Category
		if err := initializer.DB.Table("order_items oi").
			Select("c.category_name, COUNT(oi.quantity) AS quantity").
			Joins("JOIN products p ON oi.product_id = p.id").Joins("JOIN categories c ON  c.id=p.category_id").
			Group("c.category_name").
			Order("quantity DESC").
			Limit(10).
			Scan(&BestCategory).Error; err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": err,
				"code":    500,
			})
			return
		}
		for _, v := range BestCategory {
			BestList = append(BestList, gin.H{
				"categoryName": v.Category_name,
			})
		}
	}
	c.JSON(200, gin.H{
		"data":   BestList,
		"status": "Success",
	})
}
