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

// @Summary Admin Dashboard
// @Description Get admin dashboard info
// @Tags admin
// @Accept json
// @Produce json
// @Secure ApiKeyAuth
// @Success 200 {json} JSON "Welcome admin page"
// Failure 404 {json} "ErrorResponse"
// @Router /admin [get]
func AdminPage(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Welcome admin page",
	})
}

type adminDetail struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Admin Login
// @Description Authenticate admin credentials and generate JWT token for authentication
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Credentials body  adminDetail  true "Credentials for authentication ( username & password)"
// @Success 202 {json} JSON "Successfully logged"
// @Failure 401 {json} JSON "Invalid username or password"
// @Failure 501 {json} JSON "Error binding data"
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	var AdminCheck adminDetail
	var adminStore models.Admins
	err := c.Bind(&AdminCheck)
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "Error binding data",
			"code":   501,
		})
		return
	}
	if err := initializer.DB.First(&adminStore, "email=?", AdminCheck.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":     "Fail",
			"error":      "Invalid username or password",
			"code":       401,
			"error_type": "authentication_error",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminStore.Password), []byte(AdminCheck.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"status":     "Fail",
			"error":      "Invalid username or password",
			"code":       401,
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

// @Summary Admin Logout
// @Description Admin logout and clear cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {json} JSON "Logout successful"
// Failure 404 {json} JSON  "ErrorResponse"
// @Router /admin/logout [get]
func AdminLogout(c *gin.Context) {
	c.SetCookie("jwtTokenAdmin", "", -1, "", "", false, false)
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Logout Successfull",
	})
}

// @Summary Admin SignUp
// @Description Authenticated admin can create new user account
// @Tags Signup
// @Accept json
// @Produce json
// @Secure ApiKeyAuth
// @Param data body adminDetail true "Create Admin"
// @Success 200 {json} JSON "New admin created"
// Failure 404 {json} JSON  "ErrorResponse"
// @Router /admin/signup [post]
func AdminSignUp(c *gin.Context) {
	var adminSignUp models.Admins
	err := c.ShouldBindJSON(&adminSignUp)
	if err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Json binding error",
			"code":   406,
		})
		return
	}
	HashPass, err := bcrypt.GenerateFromPassword([]byte(adminSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(501, gin.H{
			"status": "Fail",
			"error":  "Hashing error",
			"code":   501,
		})
		return
	}
	adminSignUp.Password = string(HashPass)
	erro := initializer.DB.Create(&adminSignUp)
	if erro.Error != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"message": "Failed to signup",
			"code":    500,
		})
		return
	}
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "New admin added",
	})
}

// @Summary		list of users
// @Description get list of all registered admins
// @Tags	    Admin/Users
// @Accept	    json
// @Produce		json
// @Security ApiKeyAuth 
// @Success 200 {array} UpdateUserData "List of users"
// @Failure 500 {json} ErrorResponse "Failed to fetch user data"
// @Router	    /admin/user [get]
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
type  UpdateUserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Gender string  `json:"gender"`
}

// @Summary Edit user details
// @Description Edit user details based on user ID
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body UpdateUserData true "Update user info"
// @Success 200 {json} JSON "User updated successfully"
// @Failure 404 {json} JSON "User not found"
// @Failure 406 {json} JSON "Failed to bind data"
// @Failure 500 {json} JSON "Failed to update details"
// @Router /admin/user/{id} [patch]
func EditUserDetails(c *gin.Context) {
	var userEdit models.Users
	var  updateInfo UpdateUserData
	id := c.Param("ID")
	err := initializer.DB.First(&userEdit, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Can't find user",
			"code":   404,
		})
	} else {
		err := c.Bind(&updateInfo)
		if err != nil {
			c.JSON(406, gin.H{
				"status": "Fail",
				"error":  "Failed to bind data",
				"code":   406,
			})
			userEdit = models.Users{
				Name: updateInfo.Name,
				Gender: updateInfo.Gender,
				Phone: userEdit.Phone,
			}
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
// @Summary Block user
// @Description Update User Bloking status as Blocked or Unblocked
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path integer true "User ID"
// @Success 200 {json} JSON "User blocked or unblocked successfully"
// @Failure 404 {json} JSON "User not found"
// @Failure 500 {json} JSON "Server error while trying to change blocking status"
// @Router /admin/userblock/{id} [patch]
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

// @Summary Delete user
// @Description Delete an existing user from admin side
// @Tags Admin/Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path integer true "User ID"
// @Success 200 {json} JSON "User deleted successfully"
// @Failure 404 {json} JSON "User not found"
// @Failure 500 {json} JSON "Failed to delete user"
// @Router /admin/user/{id} [delete]
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
