package controller

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// =============== user checkout the items and payment ================
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
	// ============= check if given the payment method and addres =============
	paymentMethod := c.Request.FormValue("payment")
	Address, _ := strconv.ParseUint(c.Request.FormValue("address"), 10, 64)
	if paymentMethod == "" || Address == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Payment Method and Address are required",
		})
		return
	}

	// ============= stock check and amount calc ===================
	var Amount float64
	var totalAmount float64
	for _, val := range cartItems {
		discount := OfferDiscountCalc(val.ProductId)
		Amount = ((float64(val.Product.Price) - discount) * float64(val.Quantity))
		if val.Quantity > val.Product.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficent stock for product " + val.Product.Name,
			})
			return
		}
		totalAmount += Amount
	}

	// ================== coupon validation===============
	couponCode = c.Request.FormValue("coupon")
	var couponCheck models.Coupon
	var userLimitCheck models.Order
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
			totalAmount -= couponCheck.Discount
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
	// if payment method is online redirect to payment actions ===============
	if paymentMethod == "ONLINE" {
		order_id, err := PaymentHandler(orderId, int(totalAmount))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			tx.Rollback()
			return
		} else {
			c.JSON(200, gin.H{
				"Message":  "please complete the payment",
				"order id": order_id,
			})
			err := initializer.DB.Create(&models.PaymentDetails{
				Order_Id:      order_id,
				Receipt:       uint(orderId),
				PaymentStatus: "not done",
				PaymentAmount: int(totalAmount),
			}).Error
			if err != nil {
				c.JSON(200, gin.H{
					"error": "failed to store payment data",
				})
			}
		}
	}
	// ================= insert order details into databse ===================
	order := models.Order{
		Id:                 uint(orderId),
		UserId:             int(userId),
		OrderPaymentMethod: paymentMethod,
		AddressId:          int(Address),
		OrderAmount:        totalAmount,
		OrderDate:          time.Now(),
		CouponCode:         couponCode,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(500, "failed to place order")
		return
	}
	// ============ insert order items into database ==================
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
		// ============= if order is COD manage the stock ============
		if paymentMethod != "ONLINE" {
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
	}
	// =============== delete all items from user cart ==============
	if err := tx.Where("user_id =?", userId).Delete(&models.Cart{}); err.Error != nil {
		tx.Rollback()
		c.JSON(400, "faild to delete datas in cart.")
		return
	}
	//================= commit transaction whether no error ==================
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

// ============== list the orders to user ===============
func OrderView(c *gin.Context) {
	var orders []models.Order
	userId := c.GetUint("userid")
	initializer.DB.Where("user_id=?", userId).Find(&orders)
	for _, order := range orders {
		c.JSON(200, gin.H{
			"order id":       order.Id,
			"Amount":         order.OrderAmount,
			"payment method": order.OrderPaymentMethod,
			"order date":     order.OrderDate,
		})
	}
	orders = []models.Order{}
}

// ============= show the order details to user ==============
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
			"Product quantity": orderItem.Quantity,
			"Status":           orderItem.OrderStatus,
			"Address ID":       orderItem.Order.AddressId,
		})
	}
}

// ============== cancel the order if user don't want ==============
func CancelOrder(c *gin.Context) {
	var orderItem models.OrderItems
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
		// ======= update status as cancelled ======
		orderItem.OrderStatus = "cancelled"
		orderItem.OrderCancelReason = reason
		if err := tx.Save(&orderItem).Error; err != nil {
			c.JSON(500, "Failed to update status")
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
		//========== check coupon condition ============
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
			orderAmount.OrderAmount += couponRemove.Discount
			orderAmount.OrderAmount -= float64(orderItem.SubTotal)
			orderAmount.CouponCode = ""
		}
		if err := tx.Save(&orderAmount).Error; err != nil {
			c.JSON(400, gin.H{
				"Error": "failed to update order details",
			})
			tx.Rollback()
			return
		}
		var walletUpdate models.Wallet
		if err := tx.First(&walletUpdate, "user_id=?", orderAmount.UserId).Error; err != nil {
			c.JSON(501, gin.H{
				"error": "failed to fetch wallet details",
			})
			tx.Rollback()
			return
		} else {
			walletUpdate.Balance += orderAmount.OrderAmount
			tx.Save(&walletUpdate)
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(201, gin.H{
				"Message": "failed to commit transaction",
			})
			tx.Rollback()
		} else {
			c.JSON(201, gin.H{
				"Message": "Order Cancelled",
			})
		}
	}
}
