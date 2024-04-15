package controller

import (
	"fmt"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var products []models.Products

// @Summary		Landing page
// @Description	Get a list of products from the database
// @Tags			LandingPage
// @Accept			json
// @Produce		json
// @Success		200	{string}	OK
// @Router			/ [get]
func ProductsPage(c *gin.Context) {
	products = []models.Products{}
	var productList []gin.H
	err := initializer.DB.Order("products.name").Joins("Category").Find(&products).Error
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to find products",
			"code":   500,
		})
		return
	}
	for _, v := range products {
		productList = append(productList, gin.H{
			"Id":    v.ID,
			"name":  v.Name,
			"price": v.Price,
		})
	}
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   productList,
	})
}

// @Summary product details
// @Description Get a paginated list of products including product name, description, stock, price, brand name, and image.
// @Tags Products
// @Produce json
// @Param id path integer true "Product ID"
// @Success 200 {json} Suc	cessResponse
// @Failure 404 {json} JSON "Product details not found"
// @Router /product/{id} [get]
func ProductDetails(c *gin.Context) {
	var productdetails models.Products
	var quantity string
	var productDetailsShow []gin.H
	id := c.Param("ID")
	if err := initializer.DB.First(&productdetails, id).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "Can't see product",
			"code":   404,
		})
		return
	}
	if productdetails.Quantity == 0 {
		quantity = "Out of stock"
	} else {
		quantity = "Product available"
		discount := OfferDiscountCalc(int(productdetails.ID))
		rating := RatingCalc(id, c)
		var reviewView []models.Review
		if err := initializer.DB.Where("product_id=?", id).Joins("User").Find(&reviewView).Error; err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to fetch review details",
				"code":   500,
			})
			return
		}
		var relatedProducts []models.Products
		err := initializer.DB.Where("products.category_id =? AND products.id!=?", productdetails.CategoryId, id).Joins("Category").Find(&relatedProducts).Error
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "Failed to find related products",
				"code":   500,
			})
			return
		}
		productDetailsShow = append(productDetailsShow, gin.H{
			"name":        productdetails.Name,
			"price":       productdetails.Price,
			"description": productdetails.Description,
			"size":        productdetails.Size,
			"color":       productdetails.Color,
			"imageURL":    productdetails.ImagePath,
			"categoryId":  productdetails.CategoryId,
		})
		productDetailsShow = append(productDetailsShow, gin.H{"rating": rating})
		for _, v := range reviewView {
			productDetailsShow = append(productDetailsShow, gin.H{
				"userName": v.User.Name,
				"dateTime": v.Time,
				"review":   v.Review,
			})
		}
		for _, v := range relatedProducts {
			productDetailsShow = append(productDetailsShow, gin.H{
				"productName":  v.Name,
				"productPrice": v.Price,
				"productSize":  v.Size,
			})
		}
		c.JSON(200, gin.H{
			"status":                 "success",
			"data":                   productDetailsShow,
			"offer":                  discount,
			"productDiscountedPrice": productdetails.Price - discount,
			"stockStatus":            quantity,
		})

	}
}

// @Summary  Rating store
// @Description Product rating store
// @Tags Products
// @Produce  multipart/form-data
// @Param id path integer true "Product ID"
// @Param rating formData  string true "Rating value"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} JSON "Failed to create rating"
// @Router /product/rating/{id} [post]
func RatingStore(c *gin.Context) {
	var ratingValue models.Rating
	var ratingStore models.Rating
	productId := c.Param("ID")
	rating, _ := strconv.Atoi(c.Request.FormValue("rating"))
	result := initializer.DB.First(&ratingStore, "product_id=?", productId)
	if result.Error != nil {
		ratingValue.Users = 1
		ratingValue.Value = rating
		if err := initializer.DB.Create(&ratingValue).Error; err != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "failed to store data",
				"code":   400,
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Thanks for rating"})
		}
	} else {
		err := initializer.DB.Model(&ratingStore).Where("product_id=?", productId).Updates(models.Rating{
			Users: ratingStore.Users + 1,
			Value: ratingStore.Value + rating,
		})
		if err.Error != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to update data",
				"code":   500,
			})
		} else {
			c.JSON(201, gin.H{
				"status":  "Success",
				"message": "Thanks for rating",
				"code":    201,
			})
		}
	}
	ratingStore = models.Rating{}
}
func RatingCalc(id string, c *gin.Context) float64 {
	var ratingUser models.Rating
	if err := initializer.DB.First(&ratingUser, "product_id=?", id).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to get user info from database",
			"code":   500,
		})
	} else {
		averageratio := float64(ratingUser.Value) / float64(ratingUser.Users)
		ratingUser = models.Rating{}
		result := fmt.Sprintf("%.1f", averageratio)
		averageratio, _ = strconv.ParseFloat(result, 64)
		return averageratio
	}
	return 0
}

// @Summary  Review store
// @Description Product Review store
// @Tags Products
// @Produce multipart/form-data
// @Secure ApiKeyAuth
// @Param id path integer true "Product id fetch"
// @Param review formData string true "Product review"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} JSON "Failed to create rating"
// @Router /product/review/{id} [post]
func ReviewStore(c *gin.Context) {
	var reviewStore models.Review
	Id := c.Param("ID")
	ProductId, _ := strconv.Atoi(Id)
	reviewStore = models.Review{
		Review:    c.Request.FormValue("review"),
		UserId:    int(c.GetUint("userid")),
		Time:      time.Now().Format("2006-01-02"),
		ProductId: uint(ProductId),
	}
	if err := initializer.DB.Create(&reviewStore).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to store review",
			"code":   500,
		})
	} else {
		c.JSON(201, gin.H{
			"status":  "Success",
			"message": "Thank for your feedback"})
	}
}
