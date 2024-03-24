package controller

import (
	"fmt"
	"net/http"
	"project1/package/initializer"
	"project1/package/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
)

func SalesReport(c *gin.Context) {
	var sales []models.Order
	var totalamount float64
	initializer.DB.Find(&sales)
	for _, val := range sales {
		totalamount += val.OrderAmount
	}
	var salesItems []models.OrderItems
	var cancelCount int
	var totalSales int
	initializer.DB.Find(&salesItems)
	for _, val := range salesItems {
		if val.OrderStatus == "cancelled" {
			cancelCount++
		} else {
			totalSales++
		}
	}
}
func SalesReportExcel(c *gin.Context) {
	var OrderData []models.OrderItems
	if err := initializer.DB.Preload("Product").Preload("Order.User").Find(&OrderData).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to fetch sales data",
		})
		return
	}

	//============ Create a new excel file ==============
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sales Report")
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create Excel sheet",
		})
		return
	}

	// ============== Add headers to the excel sheet ==============
	headers := []string{"Order ID", "Customer Name", "Product Name", "Order Date", "Total Amount"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	//============= Add sales data ===============
	var totalAmount float32
	for _, sale := range OrderData {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(int(sale.OrderId))
		row.AddCell().Value = sale.Order.User.Name
		row.AddCell().Value = sale.Product.Name
		row.AddCell().Value = sale.Order.OrderDate.Format("2016-02-01") // Convert to string or format as needed
		row.AddCell().Value = fmt.Sprintf("%d", sale.SubTotal)          // Format total amount
		totalAmount += float32(sale.SubTotal)
	}
	// =========== total amount =============
    totalRow := sheet.AddRow()
    totalRow.AddCell()
    totalRow.AddCell()
    totalRow.AddCell()
    totalRow.AddCell().Value = "Total Amount:"
    totalRow.AddCell().Value = fmt.Sprintf("%.2f", totalAmount)

	// =============== Save the Excel file ============
	excelPath := "C:/Users/nuhma/Desktop/Week_Task/1st_project/sales_report.xlsx"
	if err := file.Save(excelPath); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save Excel file",
		})
		return
	}
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", excelPath))
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(excelPath)

	c.JSON(201,gin.H{
		"Message":"Excel file generated and sent successfully",
	})
	
}


func SalesReportPDF(c *gin.Context) {
	var OrderData []models.OrderItems
	if err := initializer.DB.Preload("Product").Preload("Order.User").Find(&OrderData).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to fetch sales data",
		})
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font and font size
	pdf.SetFont("Arial", "", 12)

	// Add headers to the PDF
	headers := []string{"Order ID", "Customer", "Product", "Order Date", "Total Amount"}
	for _, header := range headers {
		pdf.Cell(40, 10, header)
	}
	pdf.Ln(-1)

	// Add sales data to the PDF
	for _, sale := range OrderData {
		pdf.Cell(40, 10, strconv.Itoa(int(sale.OrderId)))
		pdf.Cell(40, 10, sale.Order.User.Name)
		pdf.Cell(40, 10, sale.Product.Name)
		pdf.Cell(40, 10, sale.Order.OrderDate.Format("2016-02-01"))
		pdf.Cell(40, 10, fmt.Sprintf("%d", sale.SubTotal))
		pdf.Ln(-1)
	}

	// Save the PDF file
	pdfPath := "C:/Users/nuhma/Desktop/Week_Task/1st_project/sales_report.pdf"
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF file"})
		return
	}

	// Return the PDF file as a download
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfPath))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(pdfPath)

	fmt.Println("PDF file generated and sent successfully")
}
