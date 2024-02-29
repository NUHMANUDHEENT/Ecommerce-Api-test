package routers

import (
	"project1/package/controller"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.RouterGroup) {
	//================ user authentication=======================
	r.GET("/", controller.AdminPage)
	r.POST("/login", controller.AdminLogin)

	//================User managment=======================
	r.GET("/user_managment", controller.UserList)
	r.PATCH("/user_managment/user_edit/:ID", controller.EditUserDetails)
	r.PATCH("/user_managment/user_block/:ID", controller.BlockUser)
	r.DELETE("/user_managment/user_delete/:ID", controller.DeleteUser)

	//================product managment=======================
	r.GET("/products", controller.ProductList)
	r.GET("/products/add_products", controller.AddProducts)
	r.POST("/products/add_products", controller.UploadImage)
	r.PATCH("products/edit_products/:ID", controller.EditProducts)
	r.DELETE("products/delete_products/:ID", controller.DeleteProducts)

	//================category managment=======================
	r.GET("/categories", controller.CategoryList)
	r.POST("/categories/add_category", controller.AddCategory)
	r.PATCH("/categories/edit_category/:ID", controller.EditCategories)
	r.DELETE("/categories/delete_category/:ID", controller.DeleteCategories)
	r.PATCH("/categories/block_category/:ID", controller.BlockCategory)

}
