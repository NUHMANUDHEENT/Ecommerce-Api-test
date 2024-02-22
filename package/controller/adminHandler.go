package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "welcome admin page"})
}

func AdminLogin(c *gin.Context) {
	err := c.ShouldBindJSON(&LogJs)
	if err != nil {
		c.JSON(501, gin.H{"error": "error binding data"})
	}

	if LogJs.Username == "nuhman1111" && LogJs.Password == "nuhman@1234" {
		c.JSON(202, gin.H{"message": "succesfully login"})
	} else {
		c.JSON(501, gin.H{"error": "invalid username or password"})
	}
}
func UserList(c *gin.Context) {
	var user_managment []models.Users
	initializer.DB.Order("ID").Find(&user_managment)
	for _, val := range user_managment {
		c.JSON(200, gin.H{
			"ID":         val.ID,
			"name":       val.Name,
			"username":   val.Username,
			"Email":      val.Email,
			"gender":     val.Gender,
			"created At": val.CreatedAt,
			"status":     val.Blocking,
		})
	}
}
func EditUserDetails(c *gin.Context) {
	var userEdit models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&userEdit, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		err := c.ShouldBindJSON(&userEdit)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to bindinng data"})
		} else {
			if err := initializer.DB.Save(&userEdit).Error; err != nil {
				c.JSON(500, gin.H{"error": "failed to update details"})
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
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		if blockUser.Blocking {
			blockUser.Blocking = false
			c.JSON(200, "user blocked")
		} else {
			blockUser.Blocking = true
			c.JSON(200, "user unblocked")
		}
		if err := initializer.DB.Save(&blockUser).Error; err != nil {
			c.JSON(500, "failed to block/unblock user")
		}
	}
}
func DeleteUser(c *gin.Context) {
	var deleteUser models.Users
	id := c.Param("ID")
	err := initializer.DB.First(&deleteUser, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find user"})
	} else {
		err := initializer.DB.Delete(&deleteUser).Error
		if err != nil {
			c.JSON(500, "failed to delete user")
		} else {
			c.JSON(200, "user deleted successfully")
		}
	}
}

//-------------------product managment-----------------

func ProductList(c *gin.Context) {
	var productList []models.Products
	// var checkCategory []models.Categories
	err := initializer.DB.Find(&productList).Error
	if err != nil {
		c.JSON(500, "failed to fetch details")
	} else {
		for _, val := range productList {
			c.JSON(200, gin.H{
				"Product Id":       val.ID,
				"Product Name":     val.Name,
				"Product Price":    val.Price,
				"Product Size":     val.Size,
				"Product Color":    val.Color,
				"Product Quantity": val.Quantity,
				"Category name":    val.Category.Category_name,
				// "Product Image":    val.ImagePath,
				"Product Status": val.Status,
				"category id":    val.CategoryId,
			})
		}
	}

}

func UploadImage(c *gin.Context) string {
	file, err := c.FormFile("p_imagepath")
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	imagePath := "C:/Users/nuhma/Desktop/Week_Task/1st_project/project_images/" + file.Filename
	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to upload photo"})
	}

	return imagePath
}

func AddProducts(c *gin.Context) {
	var addProduct models.Products
	var checkCategory models.Categories
	// imagepath := UploadImage(c)
	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	if err := initializer.DB.First(&checkCategory, addProduct.CategoryId).Error; err != nil {
		c.JSON(500, "no category found")
	} else {
		addProduct.Status = true
		// fmt.Println(addProduct, "________________________", imagepath)
		// addProduct.ImagePath = imagepath
		if result := initializer.DB.Create(&addProduct); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
		}
	}
}

func EditProducts(c *gin.Context) {
	var editProducts models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&editProducts, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find Product"})
	} else {
		err := c.ShouldBindJSON(&editProducts)
		if err != nil {
			c.JSON(500, "failed to bild details")
		} else {
			if err := initializer.DB.Save(&editProducts).Error; err != nil {
				c.JSON(500, "failed to edit details")
			}
			c.JSON(200, "successfully edited product")
		}
	}
}
func DeleteProducts(c *gin.Context) {
	var deleteProducts models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&deleteProducts, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find Product"})
	} else {
		err := initializer.DB.Delete(&deleteProducts).Error
		if err != nil {
			c.JSON(500, "failed to delete product")
		} else {
			c.JSON(200, "product deleted successfully")
		}
	}
}
