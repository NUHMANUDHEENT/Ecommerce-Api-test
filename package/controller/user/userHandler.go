package controller

import (
	"fmt"
	"project1/package/handler"
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var LogJs models.Users
var otp string
var RoleUser = "User"

type userDetailSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}

var userDetailsBind userDetailSignUp

// New user can signup with unique email and store details into database/
// @Summary User SignUp
// @Description SignUp new user with unique email and all other details.
// @Tags Signup
// @Accept json
// @Produce json
// @Param Credentials body userDetailSignUp true "User SignUp credentials"
// @Success 200 {json} SuccessResponse 
// @Failure 400 {json} ErrorResponse
// @Router /user/signup [post]
func UserSignUp(c *gin.Context) {
	LogJs = models.Users{}
	var otpStore models.OtpMail
	err := c.Bind(&userDetailsBind)
	if err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "json binding error",
			"code":   406,
		})
		return
	}

	if err := initializer.DB.First(&LogJs, "email=?", userDetailsBind.Email).Error; err == nil {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "Email address already exist",
			"code":   409,
		})
		return
	}
	otp = handler.GenerateOtp()
	fmt.Println("otp is ----------------", otp, "-----------------")
	err = handler.SendOtp(userDetailsBind.Email, otp)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to send otp",
			"code":   400,
		})
		return
	}
	result := initializer.DB.First(&otpStore, "email=?", userDetailsBind.Email)
	if result.Error != nil {
		otpStore = models.OtpMail{
			Otp:       otp,
			Email:     userDetailsBind.Email,
			CreatedAt: time.Now(),
			ExpireAt:  time.Now().Add(180 * time.Second),
		}
		err := initializer.DB.Create(&otpStore)
		if err.Error != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "failed to save otp details",
				"code":   400,
			})
			return
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", userDetailsBind.Email).Updates(models.OtpMail{
			Otp:      otp,
			ExpireAt: time.Now().Add(180 * time.Second),
		})
		if err.Error != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Failed to update OTP Details",
				"code":   400,
			})
			return
		}
	}
	c.JSON(202, gin.H{
		"status":  "Success",
		"message": "OTP has been sent successfully.",
		"otp":     otp,
	})
}
// After sending otp verify given otp with stored otp
// @Summary User SignUp otp verify
// @Description otp verification after given user details 
// @Tags Signup
// @Accept multipart/form-data
// @Produce json
// @Param otp formData int true "Verification otp"
// @Success 200 {json} SuccessResponse 
// @Failure 400 {json} ErrorResponse
// @Router /user/signup/otp [post]
func OtpCheck(c *gin.Context) {
	var userDataStore models.Users
	otp := c.Request.FormValue("otp")
	var existingOTP models.OtpMail
	if err := initializer.DB.Where("otp = ? AND expire_at > ?", otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Invalid or expired OTP",
			"code":   401,
		})
		return
	}
	HashPass, err := bcrypt.GenerateFromPassword([]byte(LogJs.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "hashing error",
			"code":   501,
		})
		return
	}
	userDataStore = models.Users{
		Name:     userDetailsBind.Name,
		Email:    userDetailsBind.Email,
		Password: string(HashPass),
		Phone:    userDetailsBind.Phone,
		Blocking: true,
	}
	erro := initializer.DB.Create(&userDataStore)
	if erro.Error != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  erro.Error.Error(),
			"code":   400,
		})
		return
	}
	if err := initializer.DB.Delete(&existingOTP).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "delete data failed",
			"code":   500,
		})
		return
	}
	userDataStore = models.Users{}
	if err := initializer.DB.First(&userDataStore).Error; err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "failed to fetch user details for wallet",
			"code":   501,
		})
		return
	}
	initializer.DB.Create(&models.Wallet{
		User_id: int(LogJs.ID),
	})
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "user created successfully",
	})
	userDetailsBind = userDetailSignUp{}
}
// If the otp not sended email otp resend option 
// @Summary User SignUp resend otp send
// @Description Resend otp send for signup 
// @Tags Signup
// @Produce json
// @Success 200 {json} SuccessResponse 
// @Failure 400 {json} ErrorResponse
// @Router /user/signup/resend [post]
func ResendOtp(c *gin.Context) {
	var otpStore models.OtpMail
	otp = handler.GenerateOtp()
	err := handler.SendOtp(userDetailsBind.Email, otp)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  err.Error(),
			"code":   400,
		})
	} else {
		result := initializer.DB.First(&otpStore, "email=?",userDetailsBind.Email)
		if result.Error != nil {
			otpStore = models.OtpMail{
				Otp:       otp,
				Email:     userDetailsBind.Email,
				CreatedAt: time.Now(),
				ExpireAt:  time.Now().Add(15 * time.Second),
			}
			err := initializer.DB.Create(&otpStore)
			if err.Error != nil {
				c.JSON(400, gin.H{
					"status": "fail",
					"error":  "failed to store otp",
					"code":   400})
			}
		} else {
			err := initializer.DB.Model(&otpStore).Where("email=?",userDetailsBind.Email).Updates(models.OtpMail{
				Otp:      otp,
				ExpireAt: time.Now().Add(15 * time.Second),
			})
			if err.Error != nil {
				c.JSON(400, gin.H{
					"status": "fail",
					"error":  "failed to update otp",
					"code":   400,
				})
			}
		}
	}
	c.JSON(202, gin.H{
		"status":  "success",
		"message": "OTP has been sent on your registered email id.",
	})
}

type userDetailLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates a user by checking their username and password.
// @Summary User Login
// @Description Authenticate a user by verifying their username and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Credentials body userDetailLogin true "User credentials (Username and Password)"
// @Success 200 {string} SuccessResponse "Login successful"
// @Failure 400 {string} ErrorResponse
// @Failure 401 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /user/login [post]
func UserLogin(c *gin.Context) {

	LogJs = models.Users{}
	var userPass userDetailLogin
	err := c.Bind(&userPass)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "error binding data",
			"code":   400,
		})
	}
	initializer.DB.First(&LogJs, "email=?", userPass.Username)
	err = bcrypt.CompareHashAndPassword([]byte(LogJs.Password), []byte(userPass.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "invalid username or password",
			"code":   400,
		})
		return
	} else {
		if !LogJs.Blocking {
			c.JSON(401, gin.H{
				"status":  "Success",
				"message": "User blocked"})
		} else {
			token := middleware.JwtTokenStart(c, LogJs.ID, LogJs.Email, RoleUser)
			c.SetCookie("jwtTokenUser", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "login successfully",
			})
		}
	}
}

// Authenicated user logout, cookie will be remove.
// @Summary User Logout
// @Description Authenticated user logout , and remove cookie and jwt from the client side.
// @Tags Authentication
// @Produce json
// @Success 200 {string} SuccessResponse "Login successful"
// @Failure 400 {string} ErrorResponse
// @Router /user/logout [get]
func UserLogout(c *gin.Context) {
	c.SetCookie("jwtTokenUser", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "logout Successfull",
	})
}
