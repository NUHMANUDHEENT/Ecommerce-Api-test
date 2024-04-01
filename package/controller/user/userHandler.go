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

func UserSignUp(c *gin.Context) {
	LogJs = models.Users{}
	var otpStore models.OtpMail
	err := c.ShouldBindJSON(&LogJs)
	if err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "json binding error",
			"code":   406,
		})
		return
	}

	if err := initializer.DB.First(&LogJs, "email=?", LogJs.Email).Error; err == nil {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "Email address already exist",
			"code":   409,
		})
		return
	}
	otp = handler.GenerateOtp()
	fmt.Println("----------------", otp, "-----------------")
	err = handler.SendOtp(LogJs.Email, otp)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to send otp",
			"code":   500,
		})
		return
	}
	result := initializer.DB.First(&otpStore, "email=?", LogJs.Email)
	if result.Error != nil {
		otpStore = models.OtpMail{
			Otp:       otp,
			Email:     LogJs.Email,
			CreatedAt: time.Now(),
			ExpireAt:  time.Now().Add(30 * time.Second),
		}
		err := initializer.DB.Create(&otpStore)
		if err.Error != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to save otp details",
				"code":   500,
			})
			return
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", LogJs.Email).Updates(models.OtpMail{
			Otp:      otp,
			ExpireAt: time.Now().Add(15 * time.Second),
		})
		if err.Error != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to update OTP Details",
				"code":   500,
			})
			return
		}
	}
	c.JSON(202, gin.H{
		"status":  "Success",
		"message": "OTP has been sent successfully." + otp,
	})
}

func OtpCheck(c *gin.Context) {
	var otpcheck models.OtpMail
	var otpExistTable models.OtpMail
	initializer.DB.First(&otpExistTable, "email=?", LogJs.Email)

	err := c.ShouldBindJSON(&otpcheck)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
			"error":  "failed to bind otp",
			"code":   500,
		})
	}
	var existingOTP models.OtpMail
	if err := initializer.DB.Where("otp = ? AND expire_at > ?", otpcheck.Otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Invalid or expired OTP",
			"code":   401,
		})
		return
	} else {
		fmt.Println("currect otp")
		HashPass, err := bcrypt.GenerateFromPassword([]byte(LogJs.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(501, gin.H{
				"status": "Fail",
				"error":  "hashing error",
				"code":   501,
			})
		}

		LogJs.Password = string(HashPass)
		LogJs.Blocking = true
		erro := initializer.DB.Create(&LogJs)
		if erro.Error != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  erro.Error.Error(),
				"code":   500,
			})
		} else {
			if err := initializer.DB.Delete(&otpExistTable).Error; err != nil {
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "delete data failed",
					"code":   500,
				})
			}
			if err := initializer.DB.First(&LogJs).Error; err != nil {
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
				"message": "user created successfully"})
		}
	}
}
func ResendOtp(c *gin.Context) {
	var otpStore models.OtpMail
	otp = handler.GenerateOtp()
	err := handler.SendOtp(LogJs.Email, otp)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
			"error":  err.Error(),
			"code":   500,
		})
	} else {
		result := initializer.DB.First(&otpStore, "email=?", LogJs.Email)
		if result.Error != nil {
			otpStore = models.OtpMail{
				Otp:       otp,
				Email:     LogJs.Email,
				CreatedAt: time.Now(),
				ExpireAt:  time.Now().Add(15 * time.Second),
			}
			err := initializer.DB.Create(&otpStore)
			if err.Error != nil {
				c.JSON(500, gin.H{
					"status": "fail",
					"error":  "failed to store otp",
					"code":   500})
			}
		} else {
			err := initializer.DB.Model(&otpStore).Where("email=?", LogJs.Email).Updates(models.OtpMail{
				Otp:      otp,
				ExpireAt: time.Now().Add(15 * time.Second),
			})
			if err.Error != nil {
				c.JSON(500, gin.H{
					"status": "fail",
					"error":  "failed to update otp",
					"code":   500,
				})
			}
		}
	}
	c.JSON(202, gin.H{
		"status":  "success",
		"message": "OTP has been sent on your registered email id.",
	})
}

func UserLogin(c *gin.Context) {
	LogJs = models.Users{}
	var userPass models.Users
	err := c.ShouldBindJSON(&LogJs)
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "error binding data",
			"code":   501,
		})
	}
	fmt.Println(LogJs)
	initializer.DB.First(&userPass, "email=?", LogJs.Email)
	err = bcrypt.CompareHashAndPassword([]byte(userPass.Password), []byte(LogJs.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "invalid username or password",
			"code":   501,
		})
	} else {
		if !userPass.Blocking {
			c.JSON(300, gin.H{
				"status":  "Success",
				"message": "User blocked"})
		} else {
			token := middleware.JwtTokenStart(c, userPass.ID, userPass.Email, RoleUser)
			c.SetCookie("jwtTokenUser", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "login successfully",
			})
		}
	}
}
func UserLogout(c *gin.Context) {
	c.SetCookie("jwtTokenUser", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "logout Successfull",
	})
}
