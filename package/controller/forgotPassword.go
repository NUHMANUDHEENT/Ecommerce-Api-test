package controller

import (
	"fmt"
	"net/http"
	"project1/package/handler"
	"project1/package/initializer"
	"project1/package/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// var forgotPassOtp string
var userCheck models.Users

func ForgotUserCheck(c *gin.Context) {
	userCheck = models.Users{}
	var otpStore models.OtpMail
	err := c.ShouldBindJSON(&userCheck)
	if err != nil {
		c.JSON(501, gin.H{"error": "json binding error"})
	}

	if err := initializer.DB.First(&userCheck, "email=?", userCheck.Email).Error; err != nil {
		c.JSON(501, gin.H{"error": "user not found"})
	} else {
		otp = handler.GenerateOtp()
		fmt.Println("----------------", otp, "-----------------")
		err = handler.SendOtp(userCheck.Email, otp)
		if err != nil {
			c.JSON(500, "failed to send otp")
		} else {
			c.JSON(200, "otp send to mail  "+otp)
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
					c.JSON(500, gin.H{"error": "failed to save otp details"})
				}
			} else {
				err := initializer.DB.Model(&otpStore).Where("email=?", userCheck.Email).Updates(models.OtpMail{
					Otp:      otp,
					ExpireAt: time.Now().Add(15 * time.Second),
				})
				if err.Error != nil {
					c.JSON(500, "failed too update data")
				}
			}
		}
	}
}
func ForgotOtpCheck(c *gin.Context) {
	var otpcheck models.OtpMail
	var otpExistTable models.OtpMail
	initializer.DB.First(&otpExistTable, "email=?", userCheck.Email)

	err := c.ShouldBindJSON(&otpcheck)
	if err != nil {
		c.JSON(500, "failed to bind otp details")
	}
	var existingOTP models.OtpMail
	if err := initializer.DB.Where("otp = ? AND expire_at > ?", otpcheck.Otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "Enter new password",
		})
	}
}

func NewPasswordSet(c *gin.Context) {
	var newPassSet models.Users
	err := c.ShouldBindJSON(&newPassSet)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "failed to bind data",
		})
	} else {
		HashPass, err := bcrypt.GenerateFromPassword([]byte(newPassSet.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, "failed to hash password")
		} else {
			if err := initializer.DB.Model(&userCheck).Where("email=?", userCheck.Email).Updates(models.Users{
				Password: string(HashPass),
			}).Error; err != nil {
				c.JSON(500, gin.H{
					"Error": "failed to update data",
				})
			} else {
				c.JSON(200, gin.H{
					"message": "password updated",
				})
			}
		}
	}
	userCheck = models.Users{}
}
