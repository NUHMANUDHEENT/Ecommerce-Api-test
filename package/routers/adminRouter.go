package routers

import (
	controller "project1/package/controller/admin"
	"project1/package/middleware"

	"github.com/gin-gonic/gin"
)

var roleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	//================ admin authentication=======================
	r.POST("/login", controller.AdminLogin)
	r.GET("/logout", controller.AdminLogout)
	r.POST("/signup", middleware.AuthMiddleware(roleAdmin), controller.AdminSignUp)
	r.GET("/", middleware.AuthMiddleware(roleAdmin), controller.AdminPage)

	//================User managment=======================
	r.GET("/user", controller.UserList)
	r.PATCH("/user/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditUserDetails)
	r.PATCH("/userblock/:ID", middleware.AuthMiddleware(roleAdmin), controller.BlockUser)
	r.DELETE("/user/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteUser)

	//================product managment=======================
	r.GET("/products", middleware.AuthMiddleware(roleAdmin), controller.ProductList)
	r.POST("/products", middleware.AuthMiddleware(roleAdmin), controller.AddProducts)
	r.POST("/products/image", middleware.AuthMiddleware(roleAdmin), controller.UploadImage)
	r.PATCH("products/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditProducts)
	r.DELETE("products/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteProducts)

	//================category managment=======================
	r.GET("/categories", middleware.AuthMiddleware(roleAdmin), controller.CategoryList)
	r.POST("/categories", middleware.AuthMiddleware(roleAdmin), controller.AddCategory)
	r.PATCH("/categories/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditCategories)
	r.DELETE("/categories/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteCategories)
	r.PATCH("/categories/block/:ID", middleware.AuthMiddleware(roleAdmin), controller.BlockCategory)

	//===================== Coupon managment ====================
	r.GET("/coupon", middleware.AuthMiddleware(roleAdmin), controller.CouponView)
	r.POST("/coupon", middleware.AuthMiddleware(roleAdmin), controller.CouponCreate)
	r.DELETE("/coupon/:ID", middleware.AuthMiddleware(roleAdmin), controller.CouponDelete)

	// =================== order managment ==============
	r.GET("/orders", middleware.AuthMiddleware(roleAdmin), controller.AdminOrdersView)
	r.PATCH("/orderstatus/:ID", middleware.AuthMiddleware(roleAdmin), controller.AdminOrderStatus)
	r.PATCH("/ordercancel/:ID", middleware.AuthMiddleware(roleAdmin), controller.AdminCancelOrder)

	// =================== offers management =====================
	r.GET("/offer", middleware.AuthMiddleware(roleAdmin), controller.OfferList)
	r.POST("/offer", middleware.AuthMiddleware(roleAdmin), controller.OfferAdd)
	r.DELETE("/offer/:ID", middleware.AuthMiddleware(roleAdmin), controller.OfferDelete)

	// ===================== sales report =========================
	r.GET("/sales/report", middleware.AuthMiddleware(roleAdmin), controller.SalesReport)
	r.GET("/sales/report/excel", middleware.AuthMiddleware(roleAdmin), controller.SalesReportExcel)
	r.GET("/sales/report/pdf", middleware.AuthMiddleware(roleAdmin), controller.SalesReportPDF)

	r.GET("/bestselling", controller.BestSelling)
}
