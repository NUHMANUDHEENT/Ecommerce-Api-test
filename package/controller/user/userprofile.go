package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// UserDetail returns details of the authenticated user.
// @Summary Get User Details
// @Description Get details of the authenticated user including first name, last name, username, email, phone number, and wallet balance.
// @Tags User/Profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/profile [get]
func UserProfile(c *gin.Context) {
	var userAddress []models.Address
	userId := c.GetUint("userid")
	if err := initializer.DB.Preload("User").Find(&userAddress, "user_id=?", userId).Error; err != nil {
		c.JSON(500, gin.H{
			"code":   500,
			"status": "fail",
			"error":  "failed to find address",
		})
		return
	}
	var userData []gin.H
	userData = append(userData, gin.H{
		"userDetails": gin.H{
			"name":   userAddress[0].User.Name,
			"email":  userAddress[0].User.Email,
			"gender": userAddress[0].User.Gender,
			"phone":  userAddress[0].Phone,
		},
	})
	for _, data := range userAddress {
		userData = append(userData, gin.H{
			"address": gin.H{
				"street":  data.City,
				"contry":  data.Country,
				"pincode": data.Pincode,
				"phone":   data.Phone,
			},
		})
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   userData,
	})
}

type addressUpdate struct {
	Address string `json:"user_address"`
	City    string `json:"user_city"`
	State   string `json:"user_state"`
	Pincode int    `json:"user_pincode"`
	Country string `json:"user_country"`
	Phone   int    `json:"user_phone"`
}

// User miltyple address adding to database
// @Summary Add address
// @Description  add multiple addresses for a single user
// @Tags Users/Profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param  body body addressUpdate true "New address details"
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/address/{id} [post]
func AddressStore(c *gin.Context) {
	var userCheck models.Users
	var addressBind addressUpdate
	var addressStore models.Address
	userId := c.GetUint("userid")
	if err := c.Bind(&addressBind); err != nil {
		c.JSON(400, gin.H{
			"code":   400,
			"status": "fail",
			"error":  "failed to  bind data",
		})
	}
	if err := initializer.DB.First(&userCheck, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "fail",
			"error":  "failed to find user",
			"code":   404,
		})
	} else {
		addressStore = models.Address{
			Address: addressBind.Address,
			City:    addressBind.City,
			State:   addressBind.State,
			Country: addressBind.Country,
			Pincode: addressBind.Pincode,
			Phone:   addressBind.Phone,
			UserId:  int(userId),
		}
		if result := initializer.DB.Create(&addressStore); result.Error != nil {
			c.JSON(400, gin.H{
				"code":   400,
				"status": "fail",
				"error":  "failed to find user",
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "success",
				"message": "new address added successfully",
			})
		}
	}
}

//	Edit existing address of a user by address id
//
// @Summary Edit address
// @Description This API is used for editing the existing address of a user by his address id .
// @Tags Users/Profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param  addressId path int true "Address ID"
// @Param  body body addressUpdate true "Edit address"
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/address/{id} [patch]
func AddressEdit(c *gin.Context) {
	var addressEdit models.Address
	var addressBind addressUpdate
	id := c.Param("ID")
	err := initializer.DB.First(&addressEdit, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"code":   404,
			"status": "fail",
			"error":  "failed to find address",
		})
		return
	}
	erro := c.Bind(&addressBind)
	if erro != nil {
		c.JSON(400, gin.H{
			"code":   400,
			"status": "fail",
			"error":  "failed to bind data",
		})
		return
	}
	addressEdit = models.Address{
		Address: addressBind.Address,
		City:    addressBind.City,
		State:   addressBind.State,
		Country: addressBind.Country,
		Pincode: addressBind.Pincode,
		Phone:   addressBind.Phone,
	}
	if err := initializer.DB.Save(&addressEdit).Error; err != nil {
		c.JSON(400, gin.H{
			"code":   400,
			"status": "fail",
			"error":  "failed to update",
		})
		return
	}
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "address updated successfully"})
}

// delete existing address of a user by address id
// @Summary Delete address
// @Description This API is used for Delete the existing address of a user by his address id .
// @Tags Users/Profile
// @Accept   json
// @Produce json
// @Security ApiKeyAuth
// @Param  Id path int true "Address ID"
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/address/{id} [delete]
func AddressDelete(c *gin.Context) {
	var deleteAddress models.Address
	id := c.Param("id")
	err := initializer.DB.First(&deleteAddress, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"code":   404,
			"status": "fail",
			"error":  "failed to find address",
		})
	} else {
		err := initializer.DB.Delete(&deleteAddress).Error
		if err != nil {
			c.JSON(400, gin.H{
				"code":   400,
				"status": "fail",
				"error":  "failed to delete address",
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "address deleted successfully",
			})
		}
	}
}

type userDetailUpdate struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  int    `json:"phone"`
	Gender string `json:"gender"`
}

// Edit User details or update details
// @Summary  Update profile
// @Description This API is used for updating the user's information like name , email and phone number .
// @Tags Users/Profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param addressId path int true "Address ID"
// @Param  data body userDetailUpdate true "User Details"
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/edit [patch]
func EditUserProfile(c *gin.Context) {
	var editProfile models.Users
	var userUpdate userDetailUpdate
	userId := c.GetUint("userid")
	if err := initializer.DB.First(&editProfile, userId).Error; err != nil {
		c.JSON(404, gin.H{
			"code":   404,
			"status": "fail",
			"error":  "user not found",
		})
		return
	}
	err := c.Bind(&userUpdate)
	if err != nil {
		c.JSON(400, gin.H{
			"code":   400,
			"status": "fail",
			"error":  "failed to bind data",
		})
		return
	}
	editProfile = models.Users{
		Name:   userUpdate.Name,
		Email:  userUpdate.Email,
		Phone:  userUpdate.Phone,
		Gender: userUpdate.Gender,
	}
	if err := initializer.DB.Save(&editProfile).Error; err != nil {
		c.JSON(400, gin.H{
			"code":   400,
			"status": "fail",
			"error":  "failed to update data",
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "updated data",
		})
	}
}
