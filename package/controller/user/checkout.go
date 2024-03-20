package controller

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"project1/package/handler"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	couponCode := ""
	var cartItems []models.Cart
	userId := c.GetUint("userid")
	initializer.DB.Preload("Product").Where("user_id=?", userId).Find(&cartItems)
	if len(cartItems) == 0 {
		c.JSON(404, gin.H{
			"Error": "no cart data found for this user",
		})
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
	// ================== coupon validation===============
	couponCode = c.Request.FormValue("coupon")
	var couponCheck models.Coupon
	var userLimitCheck models.Order
	// ============= stock check and amount calc ===================
	var Amount float64
	var totalAmount float64
	for _, val := range cartItems {
		Amount = (float64(val.Product.Price) * float64(val.Quantity))
		if val.Quantity > val.Product.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficent stock for product " + val.Product.Name,
			})
			return
		}
		totalAmount += Amount
	}

	if couponCode != "" {
		if err := initializer.DB.First(&userLimitCheck, "coupon_code", couponCode).Error; err == nil {
			c.JSON(200, gin.H{
				"Error": "Coupon already used",
			})
			return
		}
		if err := initializer.DB.Where(" code=? AND valid_from < ? AND valid_to > ? AND coupon_condition <= ?", couponCode, time.Now(), time.Now(), totalAmount).First(&couponCheck).Error; err != nil {
			c.JSON(200, gin.H{
				"Error": "Coupon Not valid",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"Messege": "Coupon applied",
			})
		}
	}
	// ================== order id creation =======================
	const charset = "123456789"
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err)
	}
	for i, b := range randomBytes {
		randomBytes[i] = charset[b%byte(len(charset))]
	}
	orderIdstring := string(randomBytes)
	orderId, _ := strconv.Atoi(orderIdstring)
	fmt.Println("-----", orderId)

	//================ Start the transaction ===================
	tx := initializer.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if paymentMethod == "ONLINE" {
		orderResponse, err := handler.PaymentHandler(orderId, int(totalAmount-couponCheck.Discount))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			tx.Rollback()
			return
		} else {
			c.JSON(200, gin.H{
				"Message":  "please complete the payment",
				"order id": orderResponse,
			})
		}
	}

	// handler.RazorPaymentVerification(,orderResponse)
	// //================= order details store ==================

	order := models.Order{
		Id:           uint(orderId),
		UserId:       int(userId),
		OrderPayment: paymentMethod,
		AddressId:    int(Address),
		OrderAmount:  totalAmount - couponCheck.Discount,
		OrderDate:    time.Now(),
		CouponCode:   couponCode,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(500, "failed to place order")
		return
	}
	for _, val := range cartItems {
		OrderItems := models.OrderItems{
			OrderId:     uint(orderId),
			ProductId:   val.ProductId,
			Quantity:    val.Quantity,
			SubTotal:    val.Product.Price * val.Quantity,
			OrderStatus: "pending",
		}
		if err := tx.Create(&OrderItems).Error; err != nil {
			tx.Rollback()
			c.JSON(501, gin.H{
				"Error": "failed to store items details",
			})
			return
		}
		var productQuantity models.Products
		tx.First(&productQuantity, val.ProductId)
		if err := tx.Save(val.Product).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{
				"error": "Failed to Update Product Stock",
			})
			return
		}
	}
	// if err := tx.Where("user_id =?", userId).Delete(&models.Cart{}); err.Error != nil {
	// 	tx.Rollback()
	// 	c.JSON(400, "faild to delete datas in cart.")
	// 	return
	// }
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, "failed to commit transaction")
		return
	}
	if paymentMethod != "ONLINE" {
		c.JSON(501, gin.H{
			"Order":   "Order Placed successfully",
			"Message": "Order will arrive with in 4 days",
		})
	}
}

func OrderView(c *gin.Context) {
	var orders []models.Order
	userId := c.GetUint("userid")
	initializer.DB.Preload("OrderItems").Find(&orders).Where("user_id=?", userId)
	for _, order := range orders {
		c.JSON(200, gin.H{
			"order id":       order.Id,
			"Amount":         order.OrderAmount,
			"payment method": order.OrderPayment,
			"order date":     order.OrderDate,
		})
	}
}

func OrderDetails(c *gin.Context) {
	var orderitems []models.OrderItems
	orderId := c.Param("ID")

	if err := initializer.DB.Where("order_items.order_id=?", orderId).Preload("Order").Preload("Product").Find(&orderitems).Error; err != nil {
		c.JSON(400, gin.H{
			"Error": "Can't find order details",
		})
		return
	}

	for _, orderItem := range orderitems {
		c.JSON(200, gin.H{
			"order_item Id":    orderItem.Id,
			"Product":          orderItem.ProductId,
			"Product name":     orderItem.Product.Name,
			"Order date":       orderItem.Order.OrderDate,
			"Amount":           orderItem.SubTotal,
			"Payment quantity": orderItem.Quantity,
			"Status":           orderItem.OrderStatus,
			"Address ID":       orderItem.Order.AddressId,
		})
	}
}

func CancelOrder(c *gin.Context) {
	var orderItem models.OrderItems
	var productQuantity models.Products
	orderItemId := c.Param("ID")
	reason := c.Request.FormValue("reason")
	tx := initializer.DB.Begin()
	if reason == "" {
		c.JSON(500, "please give the reason")
	} else {
		if err := tx.First(&orderItem, orderItemId).Error; err != nil {
			c.JSON(500, gin.H{
				"Error": "can't find order",
			})
			tx.Rollback()
			return
		}
		if orderItem.OrderStatus == "cancelled" {
			c.JSON(202, gin.H{
				"Message": "product already cancelled",
			})
			return
		}
		orderItem.OrderStatus = "cancelled"
		orderItem.OrderCancelReason = reason
		if err := tx.Save(&orderItem).Error; err != nil {
			c.JSON(500, "Failed to update status")
			tx.Rollback()
			return
		}
		tx.First(&productQuantity, orderItem.ProductId)
		productQuantity.Quantity += orderItem.Quantity
		if err := initializer.DB.Save(&productQuantity).Error; err != nil {
			c.JSON(500, "failed to add quantity")
			tx.Rollback()
			return
		}
		var orderAmount models.Order
		if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
			c.JSON(400, gin.H{
				"Error": "failed to find order details",
			})
			tx.Rollback()
			return
		}
		var couponRemove models.Coupon
		if orderAmount.CouponCode != "" {
			if err := initializer.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
				c.JSON(400, gin.H{
					"Error": "can't find coupon code",
				})
				tx.Rollback()
			}

		}
		if couponRemove.CouponCondition > int(orderAmount.OrderAmount) {
			newAmount := 0.0
			newAmount = float64(orderItem.SubTotal) + couponRemove.Discount
			orderAmount.OrderAmount -= newAmount
			orderAmount.CouponCode = ""
		}
		if err := tx.Save(&orderAmount).Error; err != nil {
			c.JSON(400, gin.H{
				"Error": "failed to update order details",
			})
			tx.Rollback()
			return
		}
		tx.Commit()
		c.JSON(201, gin.H{
			"Message": "Order Cancelled",
		})
	}
}
