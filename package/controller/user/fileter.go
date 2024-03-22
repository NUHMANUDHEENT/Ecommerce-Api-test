package controller

import (
	"project1/package/initializer"
	"project1/package/models"
	"strings"

	"github.com/gin-gonic/gin"
)
//============== search anda filter ==================
func SeaechProduct(c *gin.Context) {
	searchQuery := c.Query("query")
	sortBy := strings.ToLower(c.DefaultQuery("sort", "a_to_z"))

	// ======== search based query ============
	query := initializer.DB
	if searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}
 // ======== filter products given query =========
	switch sortBy {
	case "price_low_to_high":
		query = query.Order("price asc")
	case "price_high_to_low":
		query = query.Order("price desc")
	case "new_arrivals":
		query = query.Order("created_at desc")
	case "a_to_z":
		query = query.Order("name asc")
	case "z_to_a":
		query = query.Order("name desc")
	case "popularity":
		var products []models.Products
		query := `SELECT * FROM products
				JOIN (
					SELECT
						product_id,
						SUM(order_quantity) as total_quantity
					FROM
						orders
					GROUP BY
						product_id
					ORDER BY
						total_quantity DESC
					LIMIT 10
				) AS subq ON products.id = subq.product_id
				WHERE
					products.deleted_at IS NULL
				ORDER BY
					subq.total_quantity DESC`
		initializer.DB.Raw(query).Scan(&products)

		for _, v := range products {
			c.JSON(200, gin.H{
				"Name":  v.Name,
				"Price": v.Price,
				"ID":    v.ID,
			})
		}
		return
	default:
		query = query.Order("name asc")
	}
	var items []models.Products
	query.Joins("Category").Find(&items)

	for _, v := range items {
		c.JSON(200, gin.H{
			"Name":     v.Name,
			"Price":    v.Price,
			"Category": v.Category.Category_name,
			"ID":       v.ID,
		})
	}
}
