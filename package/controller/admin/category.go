package controller

import (
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	var categorylist []models.Category
	initializer.DB.Find(&categorylist)
	for _, v := range categorylist {
		c.JSON(200, gin.H{
			"Id":                   v.ID,
			"Category_name":        v.Category_name,
			"Category_description": v.Category_description,
			"Category_status":      v.Blocking,
		})
	}
}
func AddCategory(c *gin.Context) {
	var addcategory models.Category
	if err := c.ShouldBind(&addcategory); err != nil {
		c.JSON(500, gin.H{"error": "Failed to bind data"})
		return
	}
	addcategory.Blocking = true
	if result := initializer.DB.Create(&addcategory); result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to insert product"})
		return
	}
	c.JSON(200, gin.H{"message": "Product created successfully"})
}
func EditCategories(c *gin.Context) {
	var editcategory models.Category
	id := c.Param("ID")
	err := initializer.DB.First(&editcategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find category"})
		return
	}
	erro := c.ShouldBindJSON(&editcategory)
	if erro != nil {
		c.JSON(500, gin.H{
			"error": "failed to bild details"})
	} else {
		if err := initializer.DB.Save(&editcategory).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "failed to edit details"})
		}
		c.JSON(200, gin.H{
			"error": "successfully edited category"})
	}
}
func DeleteCategories(c *gin.Context) {
	var deletecategory models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&deletecategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find category"})
	} else {
		err := initializer.DB.Delete(&deletecategory).Error
		if err != nil {
			c.JSON(500, gin.H{
				"error": "failed to delete category"})
		} else {
			c.JSON(200, gin.H{
				"error": "category deleted successfully"})
		}
	}
}
func BlockCategory(c *gin.Context) {
	var blockCategory models.Category
	id := c.Param("ID")
	err := initializer.DB.First(&blockCategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "can't find Category"})
	} else {
		if blockCategory.Blocking {
			blockCategory.Blocking = false
			c.JSON(200, gin.H{
				"error": "Category blocked"})
		} else {
			blockCategory.Blocking = true
			c.JSON(200, gin.H{
				"error": "Category unblocked"})
		}
		if err := initializer.DB.Save(&blockCategory).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "failed to block/unblock Category"})
		}
	}
}
