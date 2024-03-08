package controller

import (
	"net/http"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// -------------------product managment-----------------
var AddProduct models.Products

func ProductList(c *gin.Context) {
	var productList []models.Products
	// var checkCategory []models.Categories
	err := initializer.DB.Joins("Category").Find(&productList).Error
	if err != nil {
		c.JSON(500, "failed to fetch details")
	} else {
		for _, val := range productList {
			if !val.Category.Blocking {
				continue
			} else {
				c.JSON(200, gin.H{
					"Product Id":       val.ID,
					"Product Name":     val.Name,
					"Product Price":    val.Price,
					"Product Size":     val.Size,
					"Product Color":    val.Color,
					"Product Quantity": val.Quantity,
					"Category name":    val.Category.Category_name,
					"Product Status": val.Status,
					"category id":    val.CategoryId,
				})
			}
		}
	}
}

func UploadImage(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch images"})
	}
	files := file.File["images"]
	var imagePaths []string

	for _, val := range files {
		filePath := "./images/" + val.Filename
		if err = c.SaveUploadedFile(val, filePath); err != nil {
			c.JSON(500, "faield to save images")
		}
		imagePaths = append(imagePaths, filePath)
	}
	AddProduct.ImagePath1 = imagePaths[0]
	AddProduct.ImagePath2 = imagePaths[1]
	AddProduct.ImagePath3 = imagePaths[2]
	if result := initializer.DB.Create(&AddProduct); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	}
}

func AddProducts(c *gin.Context) {
	AddProduct = models.Products{}
	var checkCategory models.Category
	if err := c.ShouldBindJSON(&AddProduct); err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	if err := initializer.DB.First(&checkCategory, AddProduct.CategoryId).Error; err != nil {
		c.JSON(500, "no category found")
	} else {
		AddProduct.Status = true
		c.JSON(http.StatusOK, gin.H{"message": "upload the images"})
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
func DeleteRecovery(c *gin.Context) {
	id := c.Param("ID")

	initializer.DB.Unscoped().Model(&models.Products{}).Where("id=?", id).Update("deleted_at", nil)
	c.JSON(200,"Recoverd.")
}