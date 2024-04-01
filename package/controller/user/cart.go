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
			"status": "Fail",
			"error":  "failed to fetch data",
			"code":   500,
		})
	} else {
		var totalDiscount = 0.0
		for _, val := range cartView {
			offerDiscount := OfferDiscountCalc(val.ProductId)
			val.Product.Price -= uint(offerDiscount)
			price := int(val.Quantity) * int(val.Product.Price)
			totalAmount += price
			totalDiscount += offerDiscount * float64(val.Quantity)
			count += 1
		}
		if totalAmount == 0 {
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "No product  in your cart.",
				"data":    nil,
				"total":   0,
			})
		} else {
			c.JSON(200, gin.H{
				"cartItems":     cartView,
				"totalProducts": count,
				"totalDiscount": totalDiscount,
				"totalAmount":   totalAmount,
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
				"status": "Fail",
				"error":  "failed to add to cart",
				"code":   500,
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "Success",
				"message": "product added to cart",
			})
		}
	} else {
		c.JSON(409, gin.H{
			"status": "Exist",
			"error":  "product already added",
			"code":   409,
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
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to fetch product stock deatails",
			"code":   500,
		})
	}

	err := initializer.DB.Where("user_id=? AND product_id=?", userId, id).First(&cartStore).Error
	if err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find product",
			"code":   404,
		})
	} else {
		cartStore.Quantity += 1
		if productStock.Quantity >= cartStore.Quantity {
			if cartStore.Quantity <= 5 {
				err := initializer.DB.Where("user_id=? AND product_id=?", userId, cartStore.ProductId).Save(&cartStore)
				if err.Error != nil {
					c.JSON(500, gin.H{
						"status": "Fail",
						"error":  "failed to add to one more",
						"code":   500,
					})
				} else {
					c.JSON(500, gin.H{
						"status":   "Success",
						"message":  "one more quantity added",
						"quantity": cartStore.Quantity,
					})
				}
			} else {
				c.JSON(500, gin.H{
					"status":   "Fail",
					"error":    "can't add more quantity ",
					"maxLimit": "You can only carry up to 5 items at a time.",
					"code":     500,
				})
			}
		} else {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "product out of stock",
				"code":   503,
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
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find product",
			"code":   404,
		})
	} else {
		cartStore.Quantity -= 1
		if cartStore.Quantity >= 1 {
			err := initializer.DB.Where("user_id=? AND product_id=?", userId, cartStore.ProductId).Save(&cartStore)
			if err.Error != nil {
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "failed to update",
					"code":   500,
				})
			} else {
				c.JSON(500, gin.H{
					"status":   "Success",
					"message":  "one more quantity removed",
					"quantity": cartStore.Quantity,
				})
			}
		} else {
			c.JSON(500, gin.H{
				"status": "Success",
				"error":  "can't remove one more",
			})
		}
	}
}

// ============== delete cart item ==============
func CartProductDelete(c *gin.Context) {
	var ProductRemove models.Cart
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.Where("product_id=? AND user_id=?", id, userId).First(&ProductRemove).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Product not added to cart",
			"code":   500,
		})
	} else {
	if err:=initializer.DB.Delete(&ProductRemove).Error;err!=nil{
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete item",
			"code":   500,
		})
		return
		}
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "product remove successfully",
		})
	}
}
