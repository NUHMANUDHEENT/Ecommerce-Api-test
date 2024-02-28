package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

var products []models.Products

func ProductsPage(c *gin.Context) {
	products = []models.Products{}
	err := initializer.DB.Joins("Category").Find(&products).Error
	if err != nil {
		c.JSON(500, "failed to fetch details")
	} else {
		for _, val := range products {
			if !val.Category.Blocking {
				continue
			} else {
				c.JSON(200, gin.H{
					"Product Id":    val.ID,
					"product Image": val.ImagePath1,
					"Product Name":  val.Name,
					"Product Price": val.Price,
				})
			}
		}
	}
}
func ProductDetails(c *gin.Context) {
	var productdetails models.Products
	id := c.Param("ID")
	if err := initializer.DB.First(&productdetails, id).Error; err != nil {
		c.JSON(500, "product not available now")
	} else {
		c.JSON(200, "product details")
		c.JSON(200, gin.H{
			"product image":         productdetails.ImagePath1,
			"product image1":        productdetails.ImagePath2,
			"product image2":        productdetails.ImagePath3,
			"Product Name":          productdetails.Name,
			"Product Size":          productdetails.Size,
			"Product Color":         productdetails.Color,
			"Product Price":         productdetails.Price,
			"Product descreiption ": productdetails.Description,
		})
		for _, val := range products {
			if productdetails.CategoryId == int(val.Category.ID) {
				c.JSON(200, "related products")
				c.JSON(200, gin.H{
					"product image": val.ImagePath1,
					"product name":  val.Name,
					"product price": val.Price,
				})
			}
		}
	}
	productdetails = models.Products{}
}
