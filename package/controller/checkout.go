package controller

import (
	"fmt"
	"net/http"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		fmt.Println("failed to parse")
	}

	var cartItems []models.Cart
	initializer.DB.Preload("Product").Where("user_id=?", id).Find(&cartItems)
	if len(cartItems) == 0 {
		c.JSON(404, "no cart data found for this user")
		return
	}

	paymentMethod := c.Request.FormValue("payment")
	Address, _ := strconv.ParseUint(c.Request.FormValue("address"), 10, 64)

	if paymentMethod == "" || Address == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Payment Method and Address are required",
		})
		return
	}

	couponCode := c.Request.FormValue("coupon")
	var couponCheck models.Coupon
	if couponCode != "" {
		if err := initializer.DB.Where(" code=? AND valid_from < ? AND valid_to > ?", couponCode, time.Now(), time.Now()).First(&couponCheck).Error; err != nil {
			c.JSON(200, gin.H{
				"Error": "Coupon Not valid",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"Messege": "Coupon applied",
			})
		}
		var totalAmount float64
		for _, val := range cartItems {
			Amount := (float64(val.Product.Price) * float64(val.Quantity))
			if Amount == 0 {
				c.JSON(500, gin.H{
					"Message": "no product found in cart",
				})
				return
			}

			if val.Quantity > val.Product.Quantity {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Insufficent stock for product " + val.Product.Name,
				})
				return
			}

			val.Product.Quantity -= val.Quantity

			order := models.Order{
				UserId:        id,
				OrderPayment:  paymentMethod,
				AddressId:     int(Address),
				ProductId:     val.ProductId,
				OrderQuantity: val.Quantity,
				OrderStatus:   "pending",
			}
			if couponCode != "" {
				Amount -= couponCheck.Discount
				order.CouponCode = couponCheck.Code
			} else {
				order.CouponCode = "not used"
			}
			order.OrderAmount = Amount
			if err := initializer.DB.Create(&order).Error; err != nil {
				c.JSON(500, "failed to place order")
				return
			}
			totalAmount += Amount

			if err := initializer.DB.Save(&val.Product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Faild to Update Product Stock",
				})
				return
			}

			if err := initializer.DB.Where("user_id =? AND product_id=?", id, val.ProductId).Delete(&models.Cart{}); err.Error != nil {
				c.JSON(http.StatusBadRequest, "faild to delete datas in cart.")
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Order Placed Successfully.",
			"Amount":  totalAmount,
		})

	}
}

func OrderView(c *gin.Context) {
	var orders []models.Order
	id := c.Param("ID")
	initializer.DB.Where("user_id=?", id).Joins("Product").Find(&orders)
	for _, order := range orders {
		c.JSON(200, gin.H{
			"ID":      order.ID,
			"Product": order.Product.Name,
			"Amount":  order.OrderAmount,
			"Status":  order.OrderStatus,
		})
	}
}

func OrderDetails(c *gin.Context) {
	var order models.Order
	id := c.Param("ID")
	initializer.DB.Preload("Product").Where("id=?", id).First(&order)
	c.JSON(200, gin.H{
		"Product":         order.Product.Name,
		"Amount":          order.OrderAmount,
		"Coupon":          order.CouponCode,
		"Status":          order.OrderStatus,
		"Payment Method":  order.OrderPayment,
		"Order Confirmed": order.Model.CreatedAt,
		"Status Updated":  order.Model.UpdatedAt,
	})
}

func CancelOrder(c *gin.Context) {
	id := c.Param("ID")
	var order models.Order
	order.OrderCancelReason = c.Request.FormValue("reason")
	if order.OrderCancelReason == "" {
		c.JSON(500, "please give the reason")
	}else{
		if err:= initializer.DB.Where("id=?", id).First(&order).Error;err!=nil{
			c.JSON(500,gin.H{
				"Error":"can't find order",
			})
			return
		}
		order.OrderStatus = "cancelled"
		initializer.DB.Save(&order)
		c.JSON(200, "Order Cancelled.")
	}
}
