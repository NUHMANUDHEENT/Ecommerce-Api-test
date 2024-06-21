package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/models"

	"github.com/gin-gonic/gin"
)

// @Summary List categories
// @Description Retrieve a list of categories from the database
// @Tags Admin/Categories
// @Accept json
// @Produce json
// @Secure ApiKeyAuth
// @Success 200 {json} JSON "List of categories"
// @Failure 400 {json} JSON "Failed to fetch category list"
// @Router /admin/categories [GET]
func CategoryList(c *gin.Context) {
	var categorylist []models.Category
	if err := initializer.DB.Find(&categorylist).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch category list",
			"code":   406,
		})
		return
	}
	var categoryShow []gin.H
	for _, v := range categorylist {
		categoryShow = append(categoryShow, gin.H{
			"id":          v.ID,
			"name":        v.Category_name,
			"description": v.Category_description,
		})
	}
	c.JSON(200, gin.H{
		"status":     "Success",
		"categories": categoryShow,
	})
}

type CategoryForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @Summary Add a new category
// @Description Add a new category to the database
// @Tags Admin/Categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CategoryForm true "Category data"
// @Success 200 {json} JSON "New category created successfully"
// @Failure 406 {json} JSON "Failed to bind data"
// @Failure 500 {json} JSON "Failed to insert category"
// @Router /admin/categories [POST]
func AddCategory(c *gin.Context) {
	var bindCategory CategoryForm
	if err := c.Bind(&bindCategory); err != nil {
		c.JSON(406, gin.H{
			"status": "Fail",
			"error":  "Failed to bind data",
			"code":   406,
		})
		return
	}

	if result := initializer.DB.Create(&models.Category{
		Category_name:        bindCategory.Name,
		Category_description: bindCategory.Description,
		Blocking:             false,
	}); result.Error != nil {
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
	})
}

// @Summary Edit a category
// @Description Edit an existing category in the database
// @Tags Admin/Categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ID path integer true "Category ID"
// @Param body body CategoryForm true "Category data for editing"
// @Success 200 {json} JSON "Successfully edited category"
// @Failure 400 {json} JSON "Invalid request format"
// @Failure 404 {json} JSON "Category not found"
// @Failure 500 {json} JSON "Failed to edit category"
// @Router /admin/categories/{ID} [patch]
func EditCategories(c *gin.Context) {
	var bindCategory CategoryForm
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
	erro := c.ShouldBindJSON(&bindCategory)
	if erro != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to bild details",
			"code":   500,
		})
		return
	}
	editcategory = models.Category{
		Category_name:        bindCategory.Name,
		Category_description: bindCategory.Description,
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
	})
}

// DeleteCategory is a function that delete a category by ID from the database and return response with status or error message
// @Summary Delete a specific category
// @Description Delete  a specific category by its ID
// @Tags Admin/Categories
// @Accept json
// @Produce   json
// @Security ApiKeyAuth
// @Param ID path integer true "The Category ID you want to delete"
// @Success 200 {json} JSON "Category deleted successfully"
// @Router /admin/categories/{ID} [delete]
func DeleteCategories(c *gin.Context) {
	var deletecategory models.Category
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

// Category blocking using a specific  category ID
// @Summary Blocking a category
// @Description Block the access of products in this category
// @Tags Admin/Categories
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Param ID path int true "The Category ID that will be blocked"
// @Success 200 {json} JSON "Category deleted successfully"
// @Failure 401 {json}  JSON "Unauthorized"
// @Router /admin/categories/block/{ID} [patch]
func BlockCategory(c *gin.Context) {
	var blockCategory models.Category
	id := c.Param("ID")
	fmt.Println("ew",id)
	err := initializer.DB.First(&blockCategory, "id=?", id)
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
			"status":  "Success",
			"message": "Category blocked"})
	} else {
		blockCategory.Blocking = true
		c.JSON(200, gin.H{
			"status": "Success",
			"error":  "Category unblocked"})
	}
	if err := initializer.DB.Save(&blockCategory).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to block/unblock Category",
			"code":   500,
		})
	}
}
