package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"
	"time"

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
		c.JSON(404, "Can't see product")
	} else {
		//prodcut full details
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
		// Stock Cheking
		if productdetails.Quantity == 0 {
			c.JSON(200, gin.H{
				"Stock Status": "Out of Stock"})
		} else {
			c.JSON(200, gin.H{
				"Stock Status": "Item is currently available"})
		}
		// rating of the prodects
		rating := RatingCalc(id, c)
		c.JSON(200, gin.H{
			"Rating": rating,
		})
		// review of the product
		ReviewView(id, c)

		// related products
		for _, val := range products {
			if productdetails.CategoryId == int(val.Category.ID) && val.ID != productdetails.ID {
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
func RatingStore(c *gin.Context) {
	var ratingValue models.Rating
	var ratingStore models.Rating
	err := c.ShouldBindJSON(&ratingValue)
	if err != nil {
		c.JSON(500, "failed to bind data")
	}
	result := initializer.DB.First(&ratingStore, "product_id=?", ratingValue.ProductId)
	if result.Error != nil {
		ratingValue.Users = 1
		if err := initializer.DB.Create(&ratingValue).Error; err != nil {
			c.JSON(500, "failed to store data")
		} else {
			c.JSON(201, "Thanks for rating")
		}
	} else {
		err := initializer.DB.Model(&ratingStore).Where("product_id=?", ratingValue.ProductId).Updates(models.Rating{
			Users: ratingStore.Users + 1,
			Value: ratingStore.Value + ratingValue.Value,
		})
		if err.Error != nil {
			c.JSON(500, "failed to update data")
		} else {
			c.JSON(201, "Thanks for rating")
		}
	}
	ratingStore = models.Rating{}
}
func RatingCalc(id string, c *gin.Context) float64 {
	var ratingUser models.Rating
	if err := initializer.DB.First(&ratingUser, "product_id=?", id).Error; err != nil {
		c.JSON(500, "failed to fetch data")
	} else {
		averageratio := float64(ratingUser.Value) / float64(ratingUser.Users)
		ratingUser = models.Rating{}
		result := fmt.Sprintf("%.1f", averageratio)
		averageratio, _ = strconv.ParseFloat(result, 64)
		return averageratio
	}
	return 0
}
func ReviewStore(c *gin.Context) {
	var reviewStore models.Review
	if err := c.ShouldBindJSON(&reviewStore); err != nil {
		c.JSON(500, "failed to bind data")
	} else {

		reviewStore.Time = time.Now().Format("2006-01-02")
		if err := initializer.DB.Create(&reviewStore).Error; err != nil {
			c.JSON(500, "failed to store review")
		} else {
			c.JSON(201, "Thank for your feedback")
		}
	}
}
func ReviewView(id string, c *gin.Context) {
	var reviewView []models.Review
	if err := initializer.DB.Joins("User").Find(&reviewView).Where("product_id=?", id).Error; err != nil {
		c.JSON(500, "failed to fetch review details")
	} else {
		productId, _ := strconv.Atoi(id)
		for _, val := range reviewView {
			if val.ProductId == uint(productId) {
				c.JSON(200, gin.H{
					"review_user": val.User.Name,
					"review_date": val.Time,
					"review":      val.Review,
				})
			}
		}
	}
}
