package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"main.go/initializer"
	"main.go/model"
)

func CreateInvoice(c *gin.Context) {
	userID := c.GetUint("userid")
	orderId := c.Param("id")
	var user model.UserModel
	if err := initializer.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "User not found",
			"code":   404,
		})
		return
	}
	var orderItem []model.OrderItems
	if err := initializer.DB.Where("order_id = ? AND order_status NOT IN (?,?)", orderId, "Cancelled", "pending").Preload("Product").Preload("Order.Address").Find(&orderItem).Error; err != nil {
		c.JSON(503, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch orders",
			"err":    err.Error(),
			"code":   503,
		})
		return
	}
	for _, order := range orderItem {
		if order.OrderStatus != "Delivered" {
			c.JSON(202, gin.H{
				"status":  "Fail",
				"message": "Order not Delivered ",
				"code":    202,
			})
			return
		}
	}
	var order model.Order
	var Discount float64
	initializer.DB.First(&order, orderId)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Ln(5)
	pdf.CellFormat(0, 0, "INVOICE", "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(30)
	pdf.Cell(10, -32, "Invoice No: "+orderId)
	pdf.Ln(5)
	pdf.Cell(10, -32, "Invoice Date: "+order.OrderDate.Format("2006-01-02"))
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(10, -32, "Bill To: ")
	pdf.Ln(5)
	pdf.Cell(10, -32, "Customer Name: "+user.Name)
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(5)
	for _, val := range orderItem {
		pdf.Cell(10, -32, "Address: "+val.Order.Address.City+", "+val.Order.Address.Street)
		pdf.Ln(5)
		pdf.Cell(10, -32, strconv.Itoa(int(val.Order.Address.Pincode)))
		pdf.Ln(5)
		pdf.Cell(10, -32, "Phone no : "+strconv.Itoa(int(user.Mobile)))
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 12)
		pdf.Ln(10)
		break
	}

	pdf.SetXY(10, 20)
	pdf.CellFormat(170, 30, "Shopify", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(12, 40, "Calicut", "", 0, "R", false, 0, "")
	pdf.CellFormat(12, 50, "2nd floor ,Ph: 0494 -24 36545", "", 0, "R", false, 0, "")
	pdf.Ln(60)

	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(20, 10, "No.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 10, "Item Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Product Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	totalAmount := 0.0
	for i, order := range orderItem {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, 10, order.Product.ProductName, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", order.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", order.Product.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", order.SubTotal), "1", 0, "R", false, 0, "")
		pdf.Ln(10)
		totalAmount += float64(order.SubTotal)
	}
	if order.ShippingCharge > 0 {
		order.OrderAmount -= float64(order.ShippingCharge)
	}
	Discount = totalAmount - order.OrderAmount
	totalAmount -= float64(Discount)
	if Discount > 0 {
		pdf.CellFormat(150, 10, "Discount:", "1", 0, "R", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%2.f", Discount), "1", 0, "R", true, 0, "")
		pdf.Ln(10)
	}
	if order.ShippingCharge > 0 {
		totalAmount += float64(order.ShippingCharge)
		pdf.CellFormat(150, 10, "Shipping charge:", "1", 0, "R", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%2.f", order.ShippingCharge), "1", 0, "R", true, 0, "")
		pdf.Ln(10)
	}
	Discount = 0
	pdf.CellFormat(150, 10, "Total Amount: ", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", totalAmount), "1", 0, "R", true, 0, "")

	dirPath := "./invoices"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to create directory",
			"err":    err.Error(),
			"code":   500,
		})
		return
	}

	pdfPath := filepath.Join(dirPath, "invoice.pdf")
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to generate PDF file",
			"err":    err.Error(),
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
	})
}
