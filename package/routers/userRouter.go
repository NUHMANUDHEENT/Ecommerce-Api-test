package routers

import (
	"net/http"
	controller "project1/package/controller/user"
	"project1/package/handler"
	"project1/package/middleware"

	"github.com/gin-gonic/gin"
)

var roleuser = "User"

func UserGroup(r *gin.RouterGroup) {
	//==============user authenticatio==============
	r.POST("/user/login", controller.UserLogin)
	r.GET("/user/logout", controller.UserLogout)
	//==================siguo
	r.POST("/user/signup", controller.UserSignUp)
	r.POST("/user/signup/otp", controller.OtpCheck)
	r.POST("/user/signup/resend", controller.ResendOtp)
	r.POST("/user/forgotpass", controller.ForgotUserCheck)
	r.POST("/user/forgotpass/otp", controller.ForgotOtpCheck)
	r.PATCH("/user/new-password", controller.NewPasswordSet)

	//============= authentication google ======================
	r.GET("/auth/login", handler.Googlelogin)
	r.GET("/auth/google/callback", handler.HandleGoogleCallback)

	// ================= product page ===============
	r.GET("/", controller.ProductsPage)
	r.GET("/product/:ID", controller.ProductDetails)
	r.POST("/product/rating/:ID", middleware.AuthMiddleware(roleuser), controller.RatingStore)
	r.POST("/product/review/:ID", middleware.AuthMiddleware(roleuser), controller.ReviewStore)

	//================user profile=================
	r.GET("/user/profile", middleware.AuthMiddleware(roleuser), controller.UserProfile)
	r.POST("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressStore)
	r.PATCH("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressEdit)
	r.DELETE("/user/address/:ID", middleware.AuthMiddleware(roleuser), controller.AddressDelete)
	r.PATCH("/user/edit", middleware.AuthMiddleware(roleuser), controller.EditUserProfile)

	//================= User cart ======================
	r.GET("/cart", middleware.AuthMiddleware(roleuser), controller.CartView)
	r.POST("/cart/:ID", middleware.AuthMiddleware(roleuser), controller.CartStore)
	r.PATCH("/cart/:ID/add", middleware.AuthMiddleware(roleuser), controller.CartProductAdd)
	r.PATCH("/cart/:ID/remove", middleware.AuthMiddleware(roleuser), controller.CartProductRemove)
	r.DELETE("/cart/:ID/delete", middleware.AuthMiddleware(roleuser), controller.CartProductDelete)

	//============================= filter products ====================
	r.GET("/filter", controller.SearchProduct)

	// =======================check out ====================
	r.POST("/checkout", middleware.AuthMiddleware(roleuser), controller.CheckOut)
	r.GET("/orders", middleware.AuthMiddleware(roleuser), controller.OrderView)
	r.GET("/orderdetails/:ID", middleware.AuthMiddleware(roleuser), controller.OrderDetails)
	r.PATCH("/ordercancel/:ID", middleware.AuthMiddleware(roleuser), controller.CancelOrder)

	//=========================== payment ==========================
	r.GET("/payment", func(c *gin.Context) {
		c.HTML(http.StatusOK, "payment.html", nil)
	})
	r.POST("/payment/confirm", controller.PaymentConfirmation)

	//=========================== wishlist =========================
	r.GET("/wishlist", middleware.AuthMiddleware(roleuser), controller.WishlistProducts)
	r.POST("/wishlist/:ID", middleware.AuthMiddleware(roleuser), controller.WishlistAdd)
	r.DELETE("/wishlist/:ID", middleware.AuthMiddleware(roleuser), controller.WishlistDelete)

	r.GET("/order/invoice/:ID", middleware.AuthMiddleware(roleuser), controller.CreateInvoice)
}
