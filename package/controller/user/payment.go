package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func PaymentPage(c *gin.Context) {
	c.HTML(200, "payment.html", nil)
}
func PaymentHandler(orderId int, amount int) (string, error) {

	client := razorpay.NewClient(os.Getenv("RAZOR_PAY_KEY"), os.Getenv("RAZOR_PAY_SECRET"))
	orderParams := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(orderId),
	}
	order, err := client.Order.Create(orderParams, nil)
	if err != nil {
		return "", errors.New("PAYMENT NOT INITIATED")
	}

	razorId, _ := order["id"].(string)
	return razorId, nil
}
func PaymentConfirmation(c *gin.Context) {
	var paymentStore models.PaymentDetails
	var paymentDetails = make(map[string]string)
	if err := c.BindJSON(&paymentDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	pd := paymentDetails

	//============== verify the signature ================
	err := RazorPaymentVerification(pd["signature"], pd["order_id"], pd["payment_id"])
	if err != nil {
		fmt.Println("-------", err)
		return
	}
	if err := initializer.DB.First(&paymentStore, "order_id=?", pd["order_id"]).Error; err != nil {
		fmt.Println("can't find payment details")
		return
	}
	paymentStore.PaymentId = pd["payment_id"]
	paymentStore.PaymentStatus = "success"
	initializer.DB.Save(&paymentStore)

	//============ quantity remove ================
	var productQuantity models.Products
	var productCheck []models.OrderItems
	if err := initializer.DB.Where("order_id=?", paymentStore.Receipt).Find(&productCheck).Error; err != nil {
		fmt.Println("cant find items")
	}
	for _, val := range productCheck {
		initializer.DB.First(&productQuantity, val.ProductId)
		productQuantity.Quantity -= val.Quantity
		if err := initializer.DB.Save(&productQuantity).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to Update Product Stock",
			})
			return
		}
	}
	fmt.Println("payment done , order placed successfully")
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := os.Getenv("RAZOR_PAY_SECRET")
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("PAYMENT FAILED")
	} else {
		return nil
	}
}
