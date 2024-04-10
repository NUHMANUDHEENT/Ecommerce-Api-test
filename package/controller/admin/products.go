package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// ================= product managment =============
var AddProduct models.Products

func ProductList(c *gin.Context) {
	var product []models.Products
	var productList []gin.H
	err := initializer.DB.Joins("Category").Find(&product).Error
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "Fail",
			"message": "failed to fetch details",
			"code":    500,
		})
		return
	}
	for _, v := range product {
		productList = append(productList, gin.H{
			"Id":    v.ID,
			"name":  v.Name,
			"price": v.Price,
		})
		c.JSON(200, gin.H{
			"status": "Success",
			"data":   productList,
		})
	}
}
func UploadImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(403, gin.H{
			"status": "Failed",
			"error":  "failed to fetch images",
			"code":   403,
		})
	}
	files := file.File["images"]
	var imagePaths []string

	for _, val := range files {
		filePath := "./images/" + val.Filename
		if err = c.SaveUploadedFile(val, filePath); err != nil {
			c.JSON(500, gin.H{
				"status":  "Fail",
				"message": "faield to save images",
				"code":    500,
			})
		}
		imagePaths = append(imagePaths, filePath)
	}
	AddProduct.ImagePath = imagePaths[0]
	// AddProduct.ImagePath2 = imagePaths[1]
	// AddProduct.ImagePath3 = imagePaths[2]
	if result := initializer.DB.Create(&AddProduct); result.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to insert product",
			"code":   500,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "Product created successfully",
			"data":    AddProduct,
		})
	}
}
func AddProducts(c *gin.Context) {
	AddProduct = models.Products{}
	var checkCategory models.Category
	if err := c.ShouldBindJSON(&AddProduct); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  err,
			"code":   500,
		})
	}
	if err := initializer.DB.First(&checkCategory, AddProduct.CategoryId).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "no category found",
			"code":   500,
		})
	} else {
		AddProduct.Status = true
		c.JSON(http.StatusOK, gin.H{
			"status":  "Continue",
			"message": "upload the images"})
	}
}
func EditProducts(c *gin.Context) {
	var editProducts models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&editProducts, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find Product",
			"code":   404,
		})
	} else {
		err := c.ShouldBindJSON(&editProducts)
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to bild details",
				"code":   500,
			})
		} else {
			if err := initializer.DB.Save(&editProducts).Error; err != nil {
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "failed to edit details",
					"code":   500,
				})
			}
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "successfully edited product",
				"data":    editProducts,
			})
		}
	}
}
func DeleteProducts(c *gin.Context) {
	var deleteProducts models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&deleteProducts, id)
	if err.Error != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "can't find Product",
			"code":   404,
		})
	} else {
		err := initializer.DB.Delete(&deleteProducts).Error
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to delete product",
				"code":   500,
			})
		} else {
			c.JSON(200, gin.H{
				"status":  "Success",
				"message": "product deleted successfully"})
		}
	}
}
func DeleteRecovery(c *gin.Context) {
	id := c.Param("ID")
	initializer.DB.Unscoped().Model(&models.Products{}).Where("id=?", id).Update("deleted_at", nil)
	c.JSON(200, "Recoverd.")
}
