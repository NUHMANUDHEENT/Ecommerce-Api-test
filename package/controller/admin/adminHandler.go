package controller

import (
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var RoleAdmin = "admin"

func AdminPage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome admin page"})
}
func AdminLogin(c *gin.Context) {
	var AdminCheck models.Admins
	var adminStore models.Admins
	err := c.ShouldBindJSON(&AdminCheck)
	if err != nil {
		c.JSON(501, gin.H{"error": "Error binding data"})
		return
	}
	if err := initializer.DB.First(&adminStore, "email=?", AdminCheck.Email).Error; err != nil {
		c.JSON(501, gin.H{"error": "Invalid username or password--"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminStore.Password), []byte(AdminCheck.Password))
	if err != nil {
		c.JSON(501, gin.H{"error": "Invalid username or password"})
		return
	}
	middleware.JwtTokenStart(c, adminStore.ID, adminStore.Email, RoleAdmin)
	c.JSON(202, gin.H{"message": "Successfully logged"})
}

func AdminLogout(c *gin.Context) {
	tokenstring := c.GetHeader("authorization")
	if tokenstring == "" {
		c.JSON(400, gin.H{
			"error": "Token not provided",
		})
		return
	}
	c.SetCookie("jwt_token", "", -1, "", "", false, false)
	c.JSON(201, gin.H{
		"message": "Logout Successfull",
	})
}
func AdminSignUp(c *gin.Context) {
	var adminSignUp models.Admins
	err := c.ShouldBindJSON(&adminSignUp)
	if err != nil {
		c.JSON(501, gin.H{"error": "Json binding error"})
		return
	}
	HashPass, err := bcrypt.GenerateFromPassword([]byte(adminSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(501, gin.H{"error": "Hashing error"})
	}
	adminSignUp.Password = string(HashPass)
	erro := initializer.DB.Create(&adminSignUp)
	if erro.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to signup",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "New admin added",
	})
}
func UserList(c *gin.Context) {
	var user_managment []models.Users
	initializer.DB.Order("ID").Find(&user_managment)
	for _, val := range user_managment {
		c.JSON(200, gin.H{
			"ID":         val.ID,
			"Name":       val.Name,
			"Email":      val.Email,
			"Gender":     val.Gender,
			"Created_At": val.CreatedAt,
			"Status":     val.Blocking,
		})
	}
}
func EditUserDetails(c *gin.Context) {
	var userEdit models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&userEdit, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "Can't find user"})
	} else {
		err := c.ShouldBindJSON(&userEdit)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to bindinng data"})
		} else {
			if err := initializer.DB.Save(&userEdit).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to update details"})
			} else {
				c.JSON(200, gin.H{"message": "User updated successfully"})
			}
		}
	}
}
func BlockUser(c *gin.Context) {
	var blockUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&blockUser, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "Can't find user"})
		return
	}
	if blockUser.Blocking {
		blockUser.Blocking = false
		c.JSON(200, gin.H{
			"message": "User blocked",
		})
	} else {
		blockUser.Blocking = true
		c.JSON(200, gin.H{
			"message": "User blocked",
		})
	}
	if err := initializer.DB.Save(&blockUser).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to Block/Unblock user",
		})
	}
}
func DeleteUser(c *gin.Context) {
	var deleteUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&deleteUser, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "Can't find user"})
		return
	}
	err = initializer.DB.Delete(&deleteUser)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete user"})
		return
	}
	c.JSON(200, gin.H{
		"message": "User deleted successfully"})
}
