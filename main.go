package main

import (
	"project1/package/initializer"
	"project1/package/routers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     
	ginSwagger "github.com/swaggo/gin-swagger" 
)

func init() {
	initializer.EnvLoad()
	initializer.LoadDatabase()
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	user := router.Group("/")
	routers.UserGroup(user)

	admin := router.Group("/admin")
	routers.AdminGroup(admin)

	router.Run(":8080")

}
