package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"time"
)

func OfferDiscountCalc(productId int) float64 {
	var OfferDiscount models.Offer
	var discount float64
	if err := initializer.DB.Where("valid_from < ? AND valid_to > ? AND product_id=?",time.Now(),time.Now(),productId).First(&OfferDiscount).Error; err == nil {
		discount = OfferDiscount.Discount
		return discount
	}
	return 0
}
