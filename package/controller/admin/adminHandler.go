package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var RoleAdmin = "Admin"

func AdminPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Welcome admin page",
	})
}
func AdminLogin(c *gin.Context) {
	var AdminCheck models.Admins
	var adminStore models.Admins
	err := c.ShouldBindJSON(&AdminCheck)
	if err != nil {
		c.JSON(501, gin.H{
			"status":    "Fail",
			"error":     "Error binding data",
			"code": 501,
		})
		return
	}
	if err := initializer.DB.First(&adminStore, "email=?", AdminCheck.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":     "Fail",
			"error":      "Invalid username or password",
			"code":  401,
			"error_type": "authentication_error",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminStore.Password), []byte(AdminCheck.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"status":     "Fail",
			"error":      "Invalid username or password",
			"code":  401,
			"error_type": "authentication_error",
		})
		return
	}
	token := middleware.JwtTokenStart(c, adminStore.ID, adminStore.Email, RoleAdmin)
	c.SetCookie("jwtTokenAdmin", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
	c.JSON(202, gin.H{
		"status":  "success",
		"message": "Successfully logged",
	})
}

func AdminLogout(c *gin.Context) {
	c.SetCookie("jwt_tokenAdmin", "", -1, "", "", false, false)
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Logout Successfull",
	})
}
func AdminSignUp(c *gin.Context) {
	var adminSignUp models.Admins
	err := c.ShouldBindJSON(&adminSignUp)
	if err != nil {
		c.JSON(406, gin.H{
			"status":    "Fail",
			"error":     "Json binding error",
			"code": 406,
		})
		return
	}
	HashPass, err := bcrypt.GenerateFromPassword([]byte(adminSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(501, gin.H{
			"status":    "Fail",
			"error":     "Hashing error",
			"code": 501,
		})
		return
	}
	adminSignUp.Password = string(HashPass)
	erro := initializer.DB.Create(&adminSignUp)
	if erro.Error != nil {
		c.JSON(500, gin.H{
			"status":    "Fail",
			"message":   "Failed to signup",
			"code": 500,
		})
		return
	}
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "New admin added",
	})
}

// @Summary Get a list of users
// @Description Get a list of users from the database
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} "OK"
// @Failure 400 {json} ErrorResponse
// @Router /admin/user [get]
func UserList(c *gin.Context) {
	var userManagment []models.Users
	err := initializer.DB.Order("ID").Find(&userManagment)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch user data",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   userManagment,
	})
}
func EditUserDetails(c *gin.Context) {
	var userEdit models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&userEdit, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Can't find user",
			"code":   404,
		})
	} else {
		err := c.ShouldBindJSON(&userEdit)
		if err != nil {
			c.JSON(406, gin.H{
				"status": "Fail",
				"error":  "Failed to bind data",
				"code":   406,
			})
		} else {
			if err := initializer.DB.Save(&userEdit).Error; err != nil {
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "Failed to update details",
					"code":   500,
				})
			} else {
				c.JSON(200, gin.H{
					"status":  "success",
					"message": "User updated successfully",
					"data":    userEdit,
				})
			}
		}
	}
}
func BlockUser(c *gin.Context) {
	var blockUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&blockUser, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Can't find user",
			"code":   404,
		})
		return
	}
	if blockUser.Blocking {
		blockUser.Blocking = false
	} else {
		blockUser.Blocking = true
	}
	if err := initializer.DB.Save(&blockUser).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Server error while trying to change blocking status of the user",
			"code":   500,
		})
		return
	}
	if blockUser.Blocking {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User blocked",
			"data":    blockUser.Blocking,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User Unblocked",
			"data":    blockUser.Blocking,
		})
	}
}
func DeleteUser(c *gin.Context) {
	var deleteUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&deleteUser, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Can't find user",
			"code":   404,
		})
		return
	}
	err = initializer.DB.Delete(&deleteUser)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete user",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "succes",
		"message": "User deleted successfully",
	})
}
