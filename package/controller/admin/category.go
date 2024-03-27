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
		"categories": categorylist,
	})
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
		c.JSON(500, gin.H{"error": "Can't find category"})
		return
	}
	erro := c.ShouldBindJSON(&editcategory)
	if erro != nil {
		c.JSON(500, gin.H{
			"error": "Failed to bild details"})
		return
	}
	if err := initializer.DB.Save(&editcategory).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to edit details"})
		return
	}
	c.JSON(200, gin.H{
		"error": "Successfully edited category"})
}
func DeleteCategories(c *gin.Context) {
	var deletecategory models.Products
	id := c.Param("ID")
	err := initializer.DB.First(&deletecategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "Can't find category"})
		return
	}
	err = initializer.DB.Delete(&deletecategory)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete category"})
	} else {
		c.JSON(200, gin.H{
			"error": "Category deleted successfully"})
	}
}
func BlockCategory(c *gin.Context) {
	var blockCategory models.Category
	id := c.Param("ID")
	err := initializer.DB.First(&blockCategory, id)
	if err.Error != nil {
		c.JSON(500, gin.H{"error": "Can't find Category"})
		return
	}
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
			"error": "Failed to block/unblock Category"})
	}
}
