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

// ForgotUserCheck is an API endpoint to send an OTP to the user's email for password recovery.
// @Summary Send OTP for password recovery
// @Description Sends an OTP to the user's email for password recovery if the email exists in the database.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "User's email address"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/forgotpass [post]
func ForgotUserCheck(c *gin.Context) {
	userCheck = models.Users{}
	var otpStore models.OtpMail
	email := c.Request.FormValue("email")

	if err := initializer.DB.First(&userCheck, "email=?", email).Error; err != nil {
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
				ExpireAt:  time.Now().Add(180 * time.Second),
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
			"message": "otp send to mail  ",
			"otp":     otp,
		})
	}
}

// ForgotOtpCheck is an API endpoint to check the validity of the OTP for password recovery.
// @Summary Check OTP validity
// @Description Checks if the provided OTP is valid and not expired for password recovery.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param otp formData string true "User's OTP"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/forgotpass/otp [post]
func ForgotOtpCheck(c *gin.Context) {
	userOTP := c.Request.FormValue("otp")

	var existingOTP models.OtpMail
	if err := initializer.DB.Where("otp = ? AND expire_at > ?", userOTP, time.Now()).First(&existingOTP).Error; err != nil {
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

// NewPasswordSet is an API endpoint to set a new password after verifying the OTP.
// @Summary Set new password
// @Description Sets a new password for the user after verifying the OTP.
// @Tags User
// @Accept  multipart/form-data
// @Produce json
// @Param password formData string true "New password"
// @Success 201 {json}  SuccessResponse
// @Failure  400 {json} ErrorResponse
// @Router /user/new-password [patch]
func NewPasswordSet(c *gin.Context) {
	password := c.Request.FormValue("password")
	if !otpValid {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "verify your email first",
			"code":   403,
		})
		return
	}
	HashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"Error":  "Failed to hash password",
			"code":   400,
		})
		return
	}
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
	userCheck = models.Users{}
	otpValid = false
}
