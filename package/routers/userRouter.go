package routers

import (
	"project1/package/controller"

	"github.com/gin-gonic/gin"
)

func UserGroup(r *gin.RouterGroup) {
	r.GET("/", controller.LoadingPage)
	r.POST("/user/signup", controller.UserSignUp)
	r.POST("/user/login",controller.UserLogin)
	r.GET("/user/signup/otp",controller.OtpCheck)
	r.POST("/user/signup/resend_otp",controller.ResendOtp)
    
}
