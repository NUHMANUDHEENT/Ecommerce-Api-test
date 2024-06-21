package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CartView retrieves and displays the items in the user's cart along with total amount and discounts.
// @Summary View cart items
// @Description Retrieves and displays the items in the user's cart along with total amount and discounts.
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} JSON Response  "successful operation"
// @Failure 400 {json} JSON  ErrorResponse "Invalid input request"
// @Router /cart [get]
func CartView(c *gin.Context) {
	var cartView []models.Cart
	var cartShow []gin.H
	userId := c.GetUint("userid")
	var totalAmount = 0.0
	var count = 0.0
	var totalDiscount = 0.0
	if err := initializer.DB.Where("user_id=?", userId).Joins("Product").Find(&cartView).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to fetch data",
			"code":   400,
		})
		return
	}
		for _, v := range cartView {
			offerDiscount := OfferDiscountCalc(v.ProductId)
			v.Product.Price -= offerDiscount
			price := int(v.Quantity) * int(v.Product.Price)
			totalAmount += float64(price)
			totalDiscount += offerDiscount * float64(v.Quantity)
			count += 1
			cartShow = append(cartShow, gin.H{
				"product": gin.H{
					"name":  v.Product.Name,
					"price": v.Product.Price,
				},
				"quantity":       v.Quantity,
				"productOffer":   offerDiscount,
				"offered amount": (v.Product.Price - offerDiscount) * float64(v.Quantity),
			})
		}
		if totalAmount == 0 {
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "No product  in your cart.",
				"data":    nil,
				"total":   0,
			})
			return
		} 
		c.JSON(200, gin.H{
			"data":     cartShow,
			"totalProducts": count,
			"totalDiscount": totalDiscount,
			"totalAmount":   totalAmount,
			"status": "Success",
		})
	}


// CartStore adds a product to the user's cart if it's not already added.
// @Summary Add product to cart
// @Description Adds a product to the user's cart if it's not already added.
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ID path string true "Product ID"
// @Success 200 {json} JSON Response "Item  was successfully added."
// @Failure 400 {json} JSON ErrorResponse  "Invalid input data."
// @Router /cart/{ID} [post]
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
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "failed to add to cart",
				"code":   400,
			})
		} else {
			c.JSON(200, gin.H{
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

// CartProductAdd increases the quantity of a product in the user's cart if it's available and within the quantity limit.
// @Summary Increase quantity of product in cart
// @Description Increases the quantity of a product in the user's cart if it's available and within the quantity limit.
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ID path string true "Product ID"
// @Success 200 {string} string "one more quantity added""
// @Failure 201 {string} string "can't add more quantity"
// @Failure 202 {string} string "product out of stock"
// @Failure 400 {string} string "failed to add to one more"
// @Failure 404 {string} string "failed to fetch product stock details/can't find product"
// @Router /cart/{ID}/add [patch]
func CartProductAdd(c *gin.Context) {
	var cartStore models.Cart
	var productStock models.Products
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.First(&productStock, id).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "failed to fetch product stock deatails",
			"code":   404,
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
					c.JSON(400, gin.H{
						"status": "Fail",
						"error":  "failed to add to one more",
						"code":   400,
					})
				} else {
					c.JSON(200, gin.H{
						"status":   "Success",
						"message":  "one more quantity added",
						"quantity": cartStore.Quantity,
					})
				}
			} else {
				c.JSON(201, gin.H{
					"status":   "Fail",
					"error":    "can't add more quantity ",
					"maxLimit": "You can only carry up to 5 items at a time.",
					"code":     201,
				})
			}
		} else {
			c.JSON(202, gin.H{
				"status": "Fail",
				"error":  "product out of stock",
				"code":   202,
			})
		}
	}
}

// CartProductRemove decreases the quantity of a product in the user's cart if it's available.
// @Summary Decrease quantity of product in cart
// @Description Decreases the quantity of a product in the user's cart if it's available.
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ID path string true "Product ID"
// @Success 200 {string} string "one more quantity removed"
// @Failure 202 {string} string "can't remove one more"
// @Failure 400 {string} string "failed to update"
// @Failure 404 {string} string "can't find product"
// @Router /cart/{ID}/remove [patch]
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
				c.JSON(400, gin.H{
					"status": "Fail",
					"error":  "failed to update",
					"code":   400,
				})
			} else {
				c.JSON(200, gin.H{
					"status":   "Success",
					"message":  "one more quantity removed",
					"quantity": cartStore.Quantity,
				})
			}
		} else {
			c.JSON(202, gin.H{
				"status": "Success",
				"error":  "can't remove one more",
			})
		}
	}
}

// CartProductDelete removes a product from the user's cart.
// @Summary Remove product from cart
// @Description Removes a product from the user's cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ID path string true "Product ID"
// @Success 200 {string}  string "Item has been deleted."
// @Failure 400 {string}  string "Failed to delete item."
// @Failure 404 {string}  string "Can't find this item in your cart."
// @Router /cart/{ID}/delete [delete]
func CartProductDelete(c *gin.Context) {
	var ProductRemove models.Cart
	userId := c.GetUint("userid")
	id := c.Param("ID")
	if err := initializer.DB.Where("product_id=? AND user_id=?", id, userId).First(&ProductRemove).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Product not added to cart",
			"code":   404,
		})
	} else {
		if err := initializer.DB.Delete(&ProductRemove).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Failed to delete item",
				"code":   400,
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "product remove successfully",
		})
	}
}
