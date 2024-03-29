package controller

import (
	"fmt"
	"net/http"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func CreateInvoice(c *gin.Context) {
	userID := c.GetUint("userid")
	orderId := c.Param("ID")
	var user models.Users
	if err := initializer.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var orderItem []models.OrderItems
	if err := initializer.DB.Where("order_id = ?", orderId).Preload("Product").Preload("Order.Address").Find(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	for _, order := range orderItem {
	if 	order.OrderStatus != "delivered"{
		c.JSON(400,gin.H{
			"message":"Order not Delivered ",
		})
	}
	}
	var order models.Order
	var Discount float64
	initializer.DB.First(&order, orderId)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 15)
	pdf.Cell(10, -32, "Invoice No: "+orderId)
	pdf.Ln(5)
	pdf.Cell(10, -32, "Invoice Date: "+"21/2022")
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(10, -32, "Bill To: ")
	pdf.SetFont("Arial", "", 15)
	pdf.Ln(5)
	pdf.Cell(10, -32, "Customer Name: "+user.Name)
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(5)
	for _, val := range orderItem {
		pdf.Cell(10, -32, "Address: "+val.Order.Address.City+", "+val.Order.Address.State)
		pdf.Ln(5)
		pdf.Cell(10, -32, strconv.Itoa(val.Order.Address.Pincode))
		pdf.Ln(5)
		pdf.Cell(10, -32, "Phone no : "+strconv.Itoa(user.Phone))
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 12)
		pdf.Ln(5)
		break
	}

	pdf.Image("./assets/logo.png", 160, 10, 30, 20, false, "", 0, "")
	pdf.SetXY(10, 20)
	pdf.CellFormat(170, 30, "Hilofy", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(12, 40, "dilka , rashka del", "", 0, "R", false, 0, "")
	pdf.CellFormat(14, 50, "15th floor ,Ph: +324 36545", "", 0, "R", false, 0, "")
	pdf.Ln(40)

	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(20, 10, "No.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "Item Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Product Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	totalAmount := 0.0
	for i, order := range orderItem {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 10, order.Product.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", order.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("$%d", order.Product.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("$%d", order.SubTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(10)
		totalAmount += float64(order.SubTotal)
	}
	Discount = totalAmount - order.OrderAmount

	if Discount > 0 {
		pdf.CellFormat(130, 10, "Coupon dicount:", "1", 0, "R", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("-$%2.f", Discount), "1", 0, "R", true, 0, "")
	}
	pdf.Ln(10)
	totalAmount -= float64(Discount)
	Discount = 0
	pdf.CellFormat(130, 10, "Total Amount: ", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("$%.2f", totalAmount), "1", 0, "R", true, 0, "")

	pdfPath := "C:/Users/nuhma/Desktop/Week_Task/1st_project/invoice.pdf"
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF file"})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfPath))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(pdfPath)

	c.JSON(200, gin.H{
		"message": "PDF file generated and sent successfully",
	})
}
