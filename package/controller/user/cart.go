package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== list the cart items =================
func CartView(c *gin.Context) {
	var cartView []models.Cart
	userId := c.GetUint("userid")
	var totalAmount = 0
	var count = 0
	if err := initializer.DB.Joins("Product").Where("user_id=?", userId).Find(&cartView).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch data",
		})
	} else {
		var totalDiscount float64
		for _, val := range cartView {
			offerDiscount := OfferDiscountCalc(val.ProductId)
			totalDiscount += offerDiscount
			c.JSON(200, gin.H{
				"offer amount":            offerDiscount,
				"product name":            val.Product.Name,
				"product image":           val.Product.ImagePath1,
				"product quantity":        val.Quantity,
				"product price":           val.Product.Price,
				"product dicounted price": val.Product.Price - uint(offerDiscount),
				"product id":              val.Product.ID,
			})
			val.Product.Price -= uint(offerDiscount)
			price := int(val.Quantity) * int(val.Product.Price)
			totalAmount += price
			totalDiscount += offerDiscount * float64(val.Quantity)
			count += 1
		}
		if totalAmount == 0 {
			c.JSON(200, gin.H{
				"Message": "No product added to cart",
			})
		} else {
			c.JSON(200, gin.H{
				"total products": count,
				"total discount": totalDiscount,
				"Total Amount":   totalAmount,
			})
		}
	}

	cartView = []models.Cart{}
}

// ============= add products into cart =============
func CartStore(c *gin.Context) {
	var cartStore models.Cart
	userId := c.GetUint("userid")
	id := c.Param("ID")
	err := initializer.DB.Where("user_id=? AND product_id=?", userId, id).First(&cartStore)
	if err.Error != nil {
		cartStore.UserId = int(userId)
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

// =============== add a specific product quantity ================
func CartProductAdd(c *gin.Context) {
	var cartStore models.Cart
	var productStock models.Products
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.First(&productStock, id).Error; err != nil {
		c.JSON(500, "failed to fetch product stock deatails")
	}

	err := initializer.DB.Where("user_id=? AND product_id=?", userId, id).First(&cartStore).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": "can't find product",
		})
	} else {
		cartStore.Quantity += 1
		if productStock.Quantity >= cartStore.Quantity {
			if cartStore.Quantity <= 5 {
				err := initializer.DB.Where("user_id=? AND product_id=?", userId, cartStore.ProductId).Save(&cartStore)
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

// ===================  remove a specific product quantity ===================
func CartProductRemove(c *gin.Context) {
	var cartStore models.Cart
	userId := c.GetUint("userid")
	id := c.Param("ID")
	err := initializer.DB.Where("user_id=? AND product_id=?", userId, id).First(&cartStore).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": "can't find product",
		})
	} else {
		cartStore.Quantity -= 1
		if cartStore.Quantity >= 1 {

			err := initializer.DB.Where("user_id=? AND product_id=?", userId, cartStore.ProductId).Save(&cartStore)
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

// ============== delete cart item ==============
func CartProductDelete(c *gin.Context) {
	var ProductRemove models.Cart
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.Where("product_id=? AND user_id=?", id, userId).Delete(&ProductRemove).Error; err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to remove product",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "product remove successfully",
		})
	}
}
