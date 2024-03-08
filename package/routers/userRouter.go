package routers

import (
	controller "project1/package/controller/user"
	"project1/package/handler"
	"project1/package/middleware"

	"github.com/gin-gonic/gin"
)

var roleuser = "user"

func UserGroup(r *gin.RouterGroup) {
	//==============user authenticatio==============
	r.GET("/user/login", controller.UserLogin)
	r.GET("/user/logout", controller.UserLogout)

	r.POST("/user/signup", controller.UserSignUp)
	r.POST("/user/signup/otp", controller.OtpCheck)
	r.POST("/user/signup/resend_otp", controller.ResendOtp)
	r.GET("/user/forgotpass", controller.ForgotUserCheck)
	r.GET("/user/forgotpass/otp", controller.ForgotOtpCheck)
	r.PATCH("/user/forgotpass", controller.NewPasswordSet)

	//============= authentication google ======================
	r.GET("/login", handler.Googlelogin)

	// ================= product page ===============
	r.GET("/", middleware.AuthMiddleware(roleuser), controller.ProductsPage)
	r.GET("/product/:ID", middleware.AuthMiddleware(roleuser), controller.ProductDetails)
	r.POST("/product/rating", middleware.AuthMiddleware(roleuser), controller.RatingStore)
	r.POST("/product/review", middleware.AuthMiddleware(roleuser), controller.ReviewStore)

	//================user profile=================
	r.GET("/user/profile/:ID", middleware.AuthMiddleware(roleuser), controller.UserProfile)
	r.POST("/user/address", middleware.AuthMiddleware(roleuser), controller.AddressStore)
	r.PATCH("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressEdit)
	r.DELETE("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressDelete)
	r.PATCH("/user/edit/:ID", middleware.AuthMiddleware(roleuser), controller.EditUserProfile)

	//================= User cart ======================
	r.POST("/cart/:ID", middleware.AuthMiddleware(roleuser), controller.CartStore)
	r.GET("/cart", middleware.AuthMiddleware(roleuser), controller.CartView)
	r.PATCH("/cart/remove/:ID", middleware.AuthMiddleware(roleuser), controller.CartProductDelete)
	r.PATCH("/cart/addquantity/:ID", controller.CartProductAdd)
	r.PATCH("/cart/removequantity/:ID", controller.CartProductRemove)

	//============================= filter ====================
	r.GET("/filter", controller.FilterPrice)

	// =======================check out ====================
	r.GET("/checkout/:ID", controller.CheckOut)
	r.GET("/orders/:ID", controller.OrderView)
	r.GET("/orderdetails/:ID", controller.OrderDetails)
	r.GET("/ordercancel/:ID", controller.CancelOrder)
}
