package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	var categorylist []models.Category
	initializer.DB.Find(&categorylist)
	c.JSON(200, gin.H{
		"status":     "Success",
		"categories": categorylist,
	})
}
func AddCategory(c *gin.Context) {
	var addcategory models.Category
	if err := c.ShouldBind(&addcategory); err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Failed to bind data",
			"code":   406,
		})
		return
	}
	addcategory.Blocking = true
	if result := initializer.DB.Create(&addcategory); result.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to insert product",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Product created successfully",
		"data":    addcategory,
	})
}
func EditCategories(c *gin.Context) {
	var editcategory models.Category
	id := c.Param("ID")
	err := initializer.DB.First(&editcategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find category",
			"code":   500,
		})
		return
	}
	erro := c.ShouldBindJSON(&editcategory)
	if erro != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to bild details",
			"code":   500,
		})
		return
	}
	if err := initializer.DB.Save(&editcategory).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to edit details",
			"code":   500,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"error":  "Successfully edited category",
		"data":   editcategory,
	})
}
func DeleteCategories(c *gin.Context) {
	var deletecategory models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&deletecategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find category",
			"code":   500,
		})
		return
	}
	err = initializer.DB.Delete(&deletecategory)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to delete category",
			"code":   500,
		})
	} else {
		c.JSON(200, gin.H{
			"status": "Success",
			"error":  "Category deleted successfully",
		})
	}
}
func BlockCategory(c *gin.Context) {
	var blockCategory models.Category
	id := c.Param("ID")
	err := initializer.DB.First(&blockCategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Can't find Category",
			"code":   500,
		})
		return
	}
	if blockCategory.Blocking {
		blockCategory.Blocking = false
		c.JSON(200, gin.H{
			"status": "Success",
			"message": "Category blocked"})
	} else {
		blockCategory.Blocking = true
		c.JSON(200, gin.H{
			"status":"Success",
			"error": "Category unblocked"})
	}
	if err := initializer.DB.Save(&blockCategory).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to block/unblock Category",
			"code":   500,
		})
	}
}
