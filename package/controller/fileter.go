package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func FilterPrice(c *gin.Context) {
	var filterStore []models.Products
	var filter string
	err := c.ShouldBind(filter)
	if err != nil {
		c.JSON(500, "failed to bind")
	}
	fmt.Println("-------", filter)
	if filter == "price" {
		if err := initializer.DB.Order("price").Find(&filterStore).Error; err != nil {
			c.JSON(500, "failed to fectch")
		}
	}
	for _, val := range filterStore {
		c.JSON(200, gin.H{
			"prodict name":  val.Name,
			"product price": val.Price,
		})
	}
}
