package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CartView(c *gin.Context) {
	var cartView []models.Cart
	var totalAmount = 0
	var count = 0
	if err := initializer.DB.Joins("Product").Where("user_id=?", middleware.UserData.ID).Find(&cartView).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch data",
		})
	} else {
		for _, val := range cartView {
			c.JSON(200, gin.H{
				"product name":     val.Product.Name,
				"product image":    val.Product.ImagePath1,
				"product quantity": val.Quantity,
				"product price":    val.Product.Price,
				"product id":       val.Product.ID,
			})
			price := int(val.Quantity) * int(val.Product.Price)
			totalAmount += price
			count += 1
		}
		if totalAmount == 0 {
			c.JSON(200, gin.H{
				"Message": "No product added to cart",
			})
		} else {
			c.JSON(200, gin.H{
				"total products": count,
				"Total Amount":   totalAmount,
			})
		}
	}

	cartView = []models.Cart{}
}
func CartStore(c *gin.Context) {
	var cartStore models.Cart
	id := c.Param("ID")
	err := initializer.DB.Where("user_id=? AND product_id=?", middleware.UserData.ID, id).First(&cartStore)
	if err.Error != nil {
		cartStore.UserId = int(middleware.UserData.ID)
		cartStore.ProductId, _ = strconv.Atoi(id)
		cartStore.Quantity = 1
		if err := initializer.DB.Create(&cartStore).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "failed to add to cart",
			})
		} else {
			c.JSON(500, gin.H{
				"message": "product added to cart",
			})
		}
	} else {
		c.JSON(500, gin.H{
			"error": "product already added",
		})
	}
}
func CartProductAdd(c *gin.Context) {
	var cartStore models.Cart
	var productStock models.Products
	id := c.Param("ID")
	if err := initializer.DB.First(&productStock, id).Error; err != nil {
		c.JSON(500, "failed to fetch product stock deatails")
	}

	err := initializer.DB.Where("user_id=? AND product_id=?", middleware.UserData.ID, id).First(&cartStore).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": "can't find product",
		})
	} else {
		cartStore.Quantity += 1
		if productStock.Quantity >= cartStore.Quantity {
			if cartStore.Quantity <= 5 {
				err := initializer.DB.Where("user_id=? AND product_id=?", middleware.UserData.ID, cartStore.ProductId).Save(&cartStore)
				if err.Error != nil {
					c.JSON(500, gin.H{
						"error": "failed to add to one more",
					})
				} else {
					c.JSON(500, gin.H{
						"quantity": cartStore.Quantity,
						"error":    "one more quantity added",
					})
				}
			} else {
				c.JSON(500, gin.H{
					"error": "can't add more quantity ",
				})
			}
		} else {
			c.JSON(500, gin.H{
				"error": "product out of stock",
			})
		}
	}
}
func CartProductRemove(c *gin.Context) {
	var cartStore models.Cart
	id := c.Param("ID")
	err := initializer.DB.Where("user_id=? AND product_id=?", middleware.UserData.ID, id).First(&cartStore).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": "can't find product",
		})
	} else {
		cartStore.Quantity -= 1
		if cartStore.Quantity >= 1 {

			err := initializer.DB.Where("user_id=? AND product_id=?", middleware.UserData.ID, cartStore.ProductId).Save(&cartStore)
			if err.Error != nil {
				c.JSON(500, gin.H{
					"error": "failed to update",
				})
			} else {
				c.JSON(500, gin.H{
					"quantity": cartStore.Quantity,
					"error":    "one more quantity removed",
				})
			}
		} else {
			c.JSON(500, gin.H{
				"error": "can't remove one more",
			})
		}
	}
}

func CartProductDelete(c *gin.Context) {
	var ProductRemove models.Cart
	id := c.Param("ID")
	if err := initializer.DB.First(&ProductRemove, "product_id=? AND user_id=?", id, middleware.UserData.ID).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "can't find product",
		})
	} else {
		fmt.Println("----------", ProductRemove)
		if err := initializer.DB.Where("product_id=? AND user_id=?", id, middleware.UserData.ID).Delete(&ProductRemove).Error; err != nil {
			c.JSON(500, gin.H{
				"Error": "failed to remove product",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "product remove successfully",
			})
		}
	}
}
