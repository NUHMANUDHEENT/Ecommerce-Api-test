package controller

import (
	"fmt"
	"project1/package/handler"
	"project1/package/initializer"
	"project1/package/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var userCheck models.Users
var otpValid = false

// =========== check user if already exist ==========
func ForgotUserCheck(c *gin.Context) {
	userCheck = models.Users{}
	var otpStore models.OtpMail
	err := c.ShouldBindJSON(&userCheck)
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "json binding error",
			"code":   501,
		})
	}
	if err := initializer.DB.First(&userCheck, "email=?", userCheck.Email).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "user not found",
			"code":   404,
		})
	} else {
		otp = handler.GenerateOtp()
		fmt.Println("----------------", otp, "-----------------")
		err = handler.SendOtp(userCheck.Email, otp)
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Otp  sending failed.",
				"code":   500,
			})
			return
		}
		result := initializer.DB.First(&otpStore, "email=?", userCheck.Email)
		if result.Error != nil {
			otpStore = models.OtpMail{
				Otp:       otp,
				Email:     userCheck.Email,
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
			}
		} else {
			if err := initializer.DB.Model(&otpStore).Where("email=?", userCheck.Email).Updates(models.OtpMail{
				Otp:      otp,
				ExpireAt: time.Now().Add(15 * time.Second),
			}).Error; err != nil {
				c.JSON(500, gin.H{
					"status": "fail",
					"error":  "failed too update data",
					"code":   500,
				})
			}
		}
		c.JSON(200, gin.H{
			"status":  "Success",
			"message": "otp send to mail  " + otp,
		})
	}
}

// VerifyOtp is used for verify
func ForgotOtpCheck(c *gin.Context) {
	var otpcheck models.OtpMail
	var otpExistTable models.OtpMail
	initializer.DB.First(&otpExistTable, "email=?", userCheck.Email)
	err := c.ShouldBindJSON(&otpcheck)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to bind otp details",
			"code":   400,
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
		otpValid = true
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Enter new password",
		})
	}
}

func NewPasswordSet(c *gin.Context) {
	if otpValid {
		var newPassSet models.Users
		err := c.ShouldBindJSON(&newPassSet)
		if err != nil {
			c.JSON(500, gin.H{
				"status": "fail",
				"error":  "failed to bind data",
				"code":   500,
			})
		} else {
			HashPass, err := bcrypt.GenerateFromPassword([]byte(newPassSet.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(500, gin.H{
					"status": "Fail",
					"Error":  "Failed to hash password",
					"code":   500,
				})
			} else {
				if err := initializer.DB.Model(&userCheck).Where("email=?", userCheck.Email).Updates(models.Users{
					Password: string(HashPass),
				}).Error; err != nil {
					c.JSON(500, gin.H{
						"status": "Fail",
						"error":  "failed to update data",
						"code":   500,
					})
				} else {
					c.JSON(201, gin.H{
						"status":  "Success",
						"message": "password updated",
					})
				}
			}
		}
		userCheck = models.Users{}
	} else {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "verify your eamil first",
			"code":   403,
		})
	}
	otpValid = false
}
