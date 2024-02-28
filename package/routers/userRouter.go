package routers

import (
	"project1/package/controller"
	"project1/package/handler"

	"github.com/gin-gonic/gin"
)

func UserGroup(r *gin.RouterGroup) {
	r.POST("/user/signup", controller.UserSignUp)
	r.POST("/user/login", controller.UserLogin)
	r.GET("/user/signup/otp", controller.OtpCheck)
	r.POST("/user/signup/resend_otp", controller.ResendOtp)

	// product page
	r.GET("/", controller.ProductsPage)
	r.GET("/product/:ID", controller.ProductDetails)
	r.POST("product/rating", controller.RatingStore)
	r.POST("product/review", controller.ReviewStore)

	// auth google
	r.GET("/login",handler.Googlelogin)
}
