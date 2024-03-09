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
	//==================siguo
	r.POST("/user/signup", controller.UserSignUp)
	r.POST("/user/signup/otp", controller.OtpCheck)
	r.POST("/user/signup/resend", controller.ResendOtp)
	r.GET("/user/forgotpass", controller.ForgotUserCheck)
	r.GET("/user/forgotpass/otp", controller.ForgotOtpCheck)
	r.PATCH("/user/forgotpass", controller.NewPasswordSet)

	//============= authentication google ======================
	r.GET("/login", handler.Googlelogin)

	// ================= product page ===============
	r.GET("/", controller.ProductsPage)
	r.GET("/product/:ID", middleware.AuthMiddleware(roleuser), controller.ProductDetails)
	r.POST("/product/rating", middleware.AuthMiddleware(roleuser), controller.RatingStore)
	r.POST("/product/review", middleware.AuthMiddleware(roleuser), controller.ReviewStore)

	//================user profile=================
	r.GET("/user/profile", middleware.AuthMiddleware(roleuser), controller.UserProfile)
	r.POST("/user/address", middleware.AuthMiddleware(roleuser), controller.AddressStore)
	r.PATCH("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressEdit)
	r.DELETE("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressDelete)
	r.PATCH("/user/edit", middleware.AuthMiddleware(roleuser), controller.EditUserProfile)

	//================= User cart ======================
	r.POST("/cart/:ID", middleware.AuthMiddleware(roleuser), controller.CartStore)
	r.GET("/cart", middleware.AuthMiddleware(roleuser), controller.CartView)
	r.PATCH("/cart/remove/:ID", middleware.AuthMiddleware(roleuser), controller.CartProductDelete)
	r.PATCH("/cart/addquantity/:ID", middleware.AuthMiddleware(roleuser), controller.CartProductAdd)
	r.PATCH("/cart/removequantity/:ID", middleware.AuthMiddleware(roleuser), controller.CartProductRemove)

	//============================= filter products ====================
	r.GET("/filter", middleware.AuthMiddleware(roleuser), controller.SeaechProduct)

	// =======================check out ====================
	r.POST("/checkout", middleware.AuthMiddleware(roleuser), controller.CheckOut)
	r.GET("/orders", middleware.AuthMiddleware(roleuser), controller.OrderView)
	r.GET("/orderdetails/:ID", middleware.AuthMiddleware(roleuser), controller.OrderDetails)
	r.PATCH("/ordercancel/:ID", middleware.AuthMiddleware(roleuser), controller.CancelOrder)
}
