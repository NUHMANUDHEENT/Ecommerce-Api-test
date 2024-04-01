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
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to find products",
			"code":   500,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   products,
	})
}

func ProductDetails(c *gin.Context) {
	var productdetails models.Products
	var quantity string
	id := c.Param("ID")
	if err := initializer.DB.First(&productdetails, id).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Find",
			"error":  "Can't see product",
			"code":   404,
		})
		return
	}
	if productdetails.Quantity == 0 {
		quantity = "Out of stock"
	} else {
		quantity = "Product available"
		discount := OfferDiscountCalc(int(productdetails.ID))
		rating := RatingCalc(id, c)
		c.JSON(200, gin.H{})
		var reviewView []models.Review
		if err := initializer.DB.Where("product_id=?", id).Joins("User").Find(&reviewView).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to fetch review details",
				"code":   500,
			})
			return
		}
		var relatedProducts []models.Products
		err := initializer.DB.Where("products.category_id =? AND products.id!=?", productdetails.CategoryId, id).Joins("Category").Find(&relatedProducts).Error
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to find related products",
				"code":   500,
			})
			return
		}
		c.JSON(200, gin.H{
			"status":                 "success",
			"data":                   productdetails,
			"offer":                  discount,
			"productDiscountedPrice": productdetails.Price - uint(discount),
			"stockStatus":            quantity,
			"rating":                 rating,
			"review":                 reviewView,
			"relatedProducts":        relatedProducts,
		})

	}
}
func RatingStore(c *gin.Context) {
	var ratingValue models.Rating
	var ratingStore models.Rating
	err := c.ShouldBindJSON(&ratingValue)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bind data",
			"code":   500,
		})
	}
	result := initializer.DB.First(&ratingStore, "product_id=?", ratingValue.ProductId)
	if result.Error != nil {
		ratingValue.Users = 1
		if err := initializer.DB.Create(&ratingValue).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to store data",
				"code":   500,
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Thanks for rating"})
		}
	} else {
		err := initializer.DB.Model(&ratingStore).Where("product_id=?", ratingValue.ProductId).Updates(models.Rating{
			Users: ratingStore.Users + 1,
			Value: ratingStore.Value + ratingValue.Value,
		})
		if err.Error != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to update data",
				"code":   500,
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Thanks for rating",
				"code":    201,
			})
		}
	}
	ratingStore = models.Rating{}
}
func RatingCalc(id string, c *gin.Context) float64 {
	var ratingUser models.Rating
	if err := initializer.DB.First(&ratingUser, "product_id=?", id).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to get user info from database",
			"code":   500,
		})
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
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bind data",
			"code":   500,
		})
	} else {
		reviewStore.Time = time.Now().Format("2006-01-02")
		reviewStore.UserId = int(c.GetUint("userid"))
		if err := initializer.DB.Create(&reviewStore).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to store review",
				"code":   500,
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Thank for your feedback"})
		}
	}
}
