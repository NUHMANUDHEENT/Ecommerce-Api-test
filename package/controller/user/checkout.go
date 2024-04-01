package controller

import (
	"crypto/rand"
	"fmt"
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
			"status":  "Fail",
			"message": "please add some items to your cart firstly.",
			"code":    404,
		})
		return
	}
	// ============= check if given payment method and addres =============
	paymentMethod := c.Request.FormValue("payment")
	Address, _ := strconv.ParseUint(c.Request.FormValue("address"), 10, 64)
	if paymentMethod == "" || Address == 0 {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Payment Method and Address are required",
			"code":   400,
		})
		return
	}
	if paymentMethod == "ONLINE" || paymentMethod == "COD" || paymentMethod == "WALLET" {
	}else{
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Give Proper Payment Method ",
			"code":   400,
		})
		return
	}

	// ============= stock check and amount calc ===================
	var Amount float64
	var totalAmount float64
	for _, val := range cartItems {
		dicount := OfferDiscountCalc(val.ProductId)
		Amount = ((float64(val.Product.Price) - dicount) * float64(val.Quantity))
		if val.Quantity > val.Product.Quantity {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Insufficent stock for product " + val.Product.Name,
				"code":   400,
			})
			return
		}
		totalAmount += Amount
	}

	// ================== coupon validation ===============
	couponCode = c.Request.FormValue("coupon")
	var couponCheck models.Coupon
	var userLimitCheck models.Order
	if couponCode != "" {
		if err := initializer.DB.First(&userLimitCheck, "coupon_code", couponCode).Error; err == nil {
			c.JSON(409, gin.H{
				"status": "Fail",
				"error":  "Coupon already used",
				"code":   409,
			})
			return
		}
		if err := initializer.DB.Where(" code=? AND valid_from < ? AND valid_to > ? AND coupon_condition <= ?", couponCode, time.Now(), time.Now(), totalAmount).First(&couponCheck).Error; err != nil {
			c.JSON(200, gin.H{
				"error": "Coupon Not valid",
			})
			return
		} else {
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
	// ============== Delivery charges ==============
	var ShippingCharge float64
	if totalAmount < 1000 {
		ShippingCharge = 40
		totalAmount += ShippingCharge
	}
	// ============== COD checking =====================
	if paymentMethod == "COD" {
		if totalAmount > 1000 {
			c.JSON(202, gin.H{
				"status":  "Fail",
				"message": "Greater than 1000 rupees should not accept COD",
				"totalAmount": totalAmount,
				"code":    202,
			})
			return
		}
	}
	// ================ wallet checking ======================
	if paymentMethod == "WALLET" {
		var walletCheck models.Wallet
		if err := initializer.DB.First(&walletCheck, "user_id=?", userId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error":  "failed to fetch wallet ",
				"code":   404,
			})
			return
		} else if walletCheck.Balance < totalAmount {
			c.JSON(202, gin.H{
				"status": "Fail",
				"error":  "insufficient balance in wallet",
				"code":   202,
			})
			return
		}

	}
	// if payment method is online redirect to payment actions ===============
	if paymentMethod == "ONLINE" {
		order_id, err := PaymentHandler(orderId, int(totalAmount))
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to create orderId",
				"code":   500,
			})
			tx.Rollback()
			return
		} else {
			c.JSON(200, gin.H{
				"status":      "Success",
				"message":     "please complete the payment",
				"totalAmount": totalAmount,
				"orderId":    order_id,
			})
			err := tx.Create(&models.PaymentDetails{
				Order_Id:      order_id,
				Receipt:       uint(orderId),
				PaymentStatus: "not done",
				PaymentAmount: int(totalAmount),
			}).Error
			if err != nil {
				c.JSON(401, gin.H{
					"status": "Fail",
					"error":  "failed to store payment data",
					"code":   401,
				})
				tx.Rollback()
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
		ShippingCharge:     ShippingCharge,
		OrderDate:          time.Now(),
		CouponCode:         couponCode,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to place order",
			"code":   500,
		})
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
				"status": "Fail",
				"error":  "failed to store items details",
				"code":   501,
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
					"status": "Fail",
					"error":  "Failed to Update Product Stock",
					"code":   500,
				})
				return
			}
		}
	}
	// =============== delete all items from user cart ==============
	if err := tx.Where("user_id =?", userId).Delete(&models.Cart{}); err.Error != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "faild to delete datas in cart.",
			"code":   400,
		})
		return
	}
	//================= commit transaction whether no error ==================
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to commit transaction",
			"code":   500,
		})
		return
	}
	if paymentMethod != "ONLINE" {
		c.JSON(501, gin.H{
			"status":      "Success",
			"Order":       "Order Placed successfully",
			"payment":     "COD",
			"totalAmount": totalAmount,
			"message":     "Order will arrive with in 4 days",
		})
	}
}

// ============== list the orders to user ===============
func OrderView(c *gin.Context) {
	var orders []models.Order
	userId := c.GetUint("userid")
	initializer.DB.Where("user_id=?", userId).Find(&orders)
	c.JSON(200, gin.H{
		"status": "success",
		"orders": orders,
	})
	orders = []models.Order{}
}

// ============= show the order details to user ==============
func OrderDetails(c *gin.Context) {
	var orderitems []models.OrderItems
	orderId := c.Param("ID")
	if err := initializer.DB.Where("order_items.order_id=?", orderId).Preload("Order").Preload("Product").Find(&orderitems).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Can't find order details",
			"code":   400,
		})
		return
	}
	var subTotal float64
	for _,val:= range  orderitems{
         subTotal +=  float64(val.SubTotal)
	}

	c.JSON(200, gin.H{
		"status":   "Success",
		"orderitems": orderitems,
		"subTotal":   subTotal,
	})
}

// ============== cancel the order if user don't want ==============
func CancelOrder(c *gin.Context) {
	var orderItem models.OrderItems
	orderItemId := c.Param("ID")
	reason := c.Request.FormValue("reason")
	tx := initializer.DB.Begin()
	if reason == "" {
		c.JSON(402, gin.H{
			"status":  "Fail",
			"message": "Please provide a valid cancellation reason.",
			"code":    402,
		})
	} else {
		if err := tx.First(&orderItem, orderItemId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error":  "can't find order",
				"code":   404,
			})
			tx.Rollback()
			return
		}
		if orderItem.OrderStatus == "cancelled" {
			c.JSON(202, gin.H{
				"status":  "Fail",
				"message": "product already cancelled",
				"code":    202,
			})
			return
		}
		// ======= update status as cancelled ======
		orderItem.OrderStatus = "cancelled"
		orderItem.OrderCancelReason = reason
		if err := tx.Save(&orderItem).Error; err != nil {
			c.JSON(500,gin.H{
				"status":"Fail",
				"error":"Failed to  save changes to database.",
				"code":500,
			})
			tx.Rollback()
			return
		}

		var orderAmount models.Order
		if err := tx.First(&orderAmount, orderItem.OrderId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error": "failed to find order details",
				"code":  404,
			})
			tx.Rollback()
			return
		}
		//========== check coupon condition ============
		var couponRemove models.Coupon
		if orderAmount.CouponCode != "" {
			if err := initializer.DB.First(&couponRemove, "code=?", orderAmount.CouponCode).Error; err != nil {
				c.JSON(404, gin.H{
					"status": "Fail",
					"error": "can't find coupon code",
					"code":   404,
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
			c.JSON(500, gin.H{
				"status": "Fail",
				"error": "failed to update order details",
				"code":   500,
			})
			tx.Rollback()
			return
		}
		var walletUpdate models.Wallet
		if err := tx.First(&walletUpdate, "user_id=?", orderAmount.UserId).Error; err != nil {
			c.JSON(501, gin.H{
				"status": "Fail",
				"error": "failed to fetch wallet details",
				"code":   501,
			})
			tx.Rollback()
			return
		} else {
			walletUpdate.Balance += orderAmount.OrderAmount
			tx.Save(&walletUpdate)
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(201, gin.H{
				"status": "Fail",
				"message": "failed to commit transaction",
				"code":    201,
			})
			tx.Rollback()
		} else {
			c.JSON(201, gin.H{
				"status":   "Success",
				"message": "Order Cancelled",
				"data":     orderItem.OrderStatus,
			})
		}
	}
}
