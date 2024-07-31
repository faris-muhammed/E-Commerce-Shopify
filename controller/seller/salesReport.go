package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
	"main.go/initializer"
	"main.go/model"
)

func SalesReport(c *gin.Context) {
	sellerId := c.GetUint("userid")

	var sales []model.Order
	var totalamount float64
	initializer.DB.Find(&sales, "seller_id=?", sellerId)
	for _, val := range sales {
		totalamount += val.OrderAmount
	}
	var salesItems []model.OrderItems
	var cancelCount int
	var totalSales int
	initializer.DB.Find(&salesItems, "seller_id=?", sellerId)
	for _, val := range salesItems {
		if val.OrderStatus == "cancelled" {
			cancelCount++
		} else {
			totalSales++
		}
	}
	c.JSON(200, gin.H{
		"TotalSalesAmount": totalamount,
		"TotalSalesCount":  totalSales,
		"TotalOrderCancel": cancelCount,
	})
}

func SalesReportExcel(c *gin.Context) {
	sellerId := c.GetUint("userid")
	var OrderData []model.OrderItems
	if err := initializer.DB.Order("").Preload("Product").Preload("Order").Find(&OrderData, "seller_id=?", sellerId).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch sales data",
			"code":   400,
		})
		return
	}

	// Create new Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sales Report")
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to create Excel sheet",
			"code":   400,
		})
		return
	}

	headers := []string{"Order ID", "Product Name", "Order Date", "Total Amount"}
	row := sheet.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.Value = header
	}

	// Add sales data
	var totalAmount float32
	for _, sale := range OrderData {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(int(sale.OrderId))
		row.AddCell().Value = sale.Product.ProductName
		row.AddCell().Value = sale.Order.OrderDate.Format("2006-01-02")
		row.AddCell().Value = fmt.Sprintf("%.2f", sale.SubTotal)
		totalAmount += float32(sale.SubTotal)
	}
	totalRow := sheet.AddRow()
	totalRow.AddCell()
	totalRow.AddCell()
	totalRow.AddCell().Value = "Total Amount:"
	totalRow.AddCell().Value = fmt.Sprintf("%.2f", totalAmount)

	// Save Excel file to local path
	dirPath := "C:/Downloads/reports"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to create directory",
			"code":   500,
		})
		return
	}

	excelPath := filepath.Join(dirPath, "sales_report.xlsx")
	if err := file.Save(excelPath); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to save Excel file",
			"code":   500,
		})
		return
	}

	// Serve the file for download
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(excelPath)))
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(excelPath)

	// Return JSON response
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "Excel file generated and sent successfully",
		"code":    201,
	})
}

func SalesReportPDF(c *gin.Context) {
	sellerId := c.GetUint("userid")
	var OrderData []model.OrderItems
	if err := initializer.DB.Preload("Product").Preload("Order").Find(&OrderData, "seller_id=?", sellerId).Error; err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch sales data",
			"code":   500,
		})
		return
	}
	// ======= create new pdf doc =========
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	headers := []string{"Order ID", "Product", "Order Date", "Total Amount"}
	for _, header := range headers {
		pdf.Cell(50, 10, header)
	}
	pdf.Ln(-1)

	// ========== add sales data ===========
	for _, sale := range OrderData {
		pdf.Cell(50, 10, strconv.Itoa(int(sale.OrderId)))
		pdf.Cell(50, 10, sale.Product.ProductName)
		pdf.Cell(50, 10, sale.Order.OrderDate.Format("2006-01-02"))
		pdf.Cell(50, 10, fmt.Sprintf("%.2f", sale.SubTotal))
		pdf.Ln(-1)
	}

	// ============== save doc into local ================
	dirPath := "C:/Downloads/reports"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to create directory",
			"code":   500,
		})
		return
	}

	pdfPath := filepath.Join(dirPath, "sales_report.pdf")
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to generate PDF file",
			"code":   500,
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfPath))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(pdfPath)

	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "PDF file generated and sent successfully",
		"code":    200,
	})
}
