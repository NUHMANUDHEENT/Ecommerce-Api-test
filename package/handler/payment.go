package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func PaymentPage(c *gin.Context) {
	c.HTML(200, "payment.html", nil)
}
func PaymentHandler(orderId int, amount int) (string, error) {
	client := razorpay.NewClient("rzp_test_CuoceDh2rleHvj", "4hciN7tsy9W3dcRV6DK3gZcN")
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

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := "YOUR_SECRET"
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
