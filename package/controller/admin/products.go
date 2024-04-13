package controller

import (
	"net/http"
	controller "project1/package/controller/user"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Product list  all products
// @Summary   List all products
// @Description  get a list of all products from the database
// @Tags Admin/Products
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} JSON Response "The order status has been changed successfully "
// @Failure 400 {json} JSON  ErrorResponse "An error occurred while updating status of the order."
// @Router /admin/products [get]
func ProductList(c *gin.Context) {
	var product []models.Products
	var productList []gin.H
	err := initializer.DB.Joins("Category").Find(&product).Error
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "Fail",
			"message": "failed to fetch details",
			"code":    400,
		})
		return
	}
	for _, v := range product {
		discount := controller.OfferDiscountCalc(int(v.ID))
		productList = append(productList, gin.H{
			"Id":       v.ID,
			"name":     v.Name,
			"price":    v.Price - discount,
			"quantity": v.Quantity,
			"category": v.CategoryId,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   productList,
	})
}

// AddProducts adds a new product with images.
// @Summary Add a new product with images
// @Description Adds a new product with images and other details
// @Tags Admin/Products
// @Accept multipart/form-data
// @Security  ApiKeyAuth
// @Param name formData string true "Product name"
// @Param price formData integer true "Product price"
// @Param size formData string true "Product size"
// @Param color formData string true "Product color"
// @Param quantity formData integer true "Product quantity"
// @Param description formData string true "Product description"
// @Param categoryId formData  int true "Category ID of the product"
// @Param images formData []file true "Product images"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/product [post]
func AddProducts(c *gin.Context) {
	var AddProduct models.Products
	var checkCategory models.Category
	File, err := c.MultipartForm()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  err,
			"code":   500,
		})
		return
	}
	AddProduct.CategoryId, _ = strconv.Atoi(c.Request.FormValue("categoryId"))
	if err := initializer.DB.First(&checkCategory, AddProduct.CategoryId).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "no category found",
			"code":   400,
		})
		return
	}
	AddProduct.Name = c.Request.FormValue("name")
	AddProduct.Price, _ = strconv.ParseFloat(c.Request.FormValue("price"), 64)
	AddProduct.Size = c.Request.FormValue("size")
	AddProduct.Color = c.Request.FormValue("color")
	Quantity, _ := strconv.Atoi(c.Request.FormValue("quantity"))
	AddProduct.Quantity = uint(Quantity)
	AddProduct.Description = c.Request.FormValue("description")
	AddProduct.Status = true
	images := File.File["images"]
	for _, v := range images {
		filePath := "./assets/" + v.Filename
		if err = c.SaveUploadedFile(v, filePath); err != nil {
			c.JSON(400, gin.H{
				"status":  "Fail",
				"message": "faield to save images",
				"code":    400,
			})
		}
		AddProduct.ImagePath = append(AddProduct.ImagePath, filePath)
	}
	if err := initializer.DB.Create(&AddProduct); err.Error != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to store product",
			"code":   400,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Product created successfully",
	})
}

// Edit existing product usign  ID as parameter
// @Summary  edit a product by id
// @Description Edit  an existing product using its unique identifier
// @Tags Admin/Products
// @Accept multipart/form-data
// @Produce json
// @Security  ApiKeyAuth
// @Param id  path int true "product Id"
// @Param name formData string true "Product name"
// @Param price formData integer true "Product price"
// @Param size formData string true "Product size"
// @Param color formData string true "Product color"
// @Param quantity formData integer true "Product quantity"
// @Param description formData string true "Product description"
// @Param categoryId formData  int true "Category ID of the product"
// @Param images formData []file true "Product images"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/product/{id} [patch]
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
		return
	}
	file, erro := c.MultipartForm()
	if erro != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to bild details",
			"code":   500,
		})
		return
	}
	editProducts.Name = c.Request.FormValue("name")
	editProducts.Price, _ = strconv.ParseFloat(c.Request.FormValue("price"), 64)
	editProducts.Size = c.Request.FormValue("size")
	editProducts.Color = c.Request.FormValue("color")
	Quantity, _ := strconv.Atoi(c.Request.FormValue("quantity"))
	editProducts.Quantity = uint(Quantity)
	editProducts.Description = c.Request.FormValue("description")
	images := file.File["image"]
	for _, img := range images {
		filePath := "./assets/" + img.Filename
		if err := c.SaveUploadedFile(img, filePath); err != nil {
			c.JSON(400, gin.H{
				"status":  "Fail",
				"message": "faield to save images",
				"code":    400,
			})
		}
		editProducts.ImagePath = append(editProducts.ImagePath, filePath)
	}
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
		"data":    "",
	})
}

// Delete Product is a function that delete the specific product by id from database
// @Summary   delete a product by its ID
// @Description  delete a product by its ID
// @Tags Admin/Products
// @Produce  json
// @Security  ApiKeyAuth
// @Param  id path int true "product's id"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/product/{id} [delete]
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
