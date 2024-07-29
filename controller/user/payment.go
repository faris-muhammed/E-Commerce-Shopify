package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"main.go/initializer"
	"main.go/model"
)

func RazorPay(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------RAZOR PAY------------------------")

	var payment model.PaymentDetails
	var orderitems []model.OrderItems
	var detail string

	orderId := c.Query("id")
	// Logged := c.GetUint("userid")

	if err := initializer.DB.Preload("User").First(&payment, "payment_id=?", orderId).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Order not found!",
			"Data":    gin.H{},
		})
		return
	}
	// if payment.UserId != Logged {
	// 	c.JSON(404, gin.H{
	// 		"Status":  "Error!",
	// 		"Code":    404,
	// 		"Message": "Order not found!",
	// 		"Data":    gin.H{},
	// 	})
	// 	return
	// }
	if err := initializer.DB.Preload("Product").Find(&orderitems, "order_id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Message": "Cannot fetch Order Items!",
			"Data":    gin.H{},
		})
		return
	}
	for i := 0; i < len(orderitems); i++ {
		if i == (len(orderitems) - 1) {
			detail += "and " + orderitems[i].Product.ProductName
		} else {
			detail += orderitems[i].Product.ProductName + ", "
		}
	}

	c.HTML(200, "payment.html", gin.H{
		"Order":  orderId,
		"Amount": payment.PaymentAmount,
		"Key":    os.Getenv("RAZOR_PAY_KEY"),
		"Name":   payment.User.Name,
		"Eamil":  payment.User.Email,
		// "Phone":   payment.User.Phone,
		"Product": "Your products " + detail + ". Pay for them now!",
	})
}

func PaymentHandler(orderId int, amount float64) (string, error) {

	client := razorpay.NewClient(os.Getenv("RAZOR_PAY_KEY"), os.Getenv("RAZOR_PAY_SECRET"))
	orderParams := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(orderId),
	}
	order, err := client.Order.Create(orderParams, nil)
	if err != nil {
		return "", errors.New("PAYMENT NOT INITIATED")
	}

	razorId, _ := order["id"].(string)
	return razorId, nil
}

type Razor struct {
	Order     string `json:"OrderID"`
	Payment   string `json:"PaymentID"`
	Signature string `json:"Signature"`
	Status    string `json:"Status"`
}

func RazorPayVerify(c *gin.Context) {
	fmt.Println("")
	fmt.Println("-----------------------------PAYMENT VERIFY------------------------")

	var verify Razor
	var order []model.OrderItems
	var payment model.PaymentDetails
	var productQuantity model.ProductDetails

	// Logged := c.GetUint("userid")

	err := c.ShouldBindJSON(&verify)
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Couldn't bind JSON data!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't bind JSON data!")
		fmt.Println("Received JSON:", verify)
		return
	}

	er := initializer.DB.First(&payment, "payment_id=?", verify.Order).Error
	if er != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   er.Error(),
			"Message": "No such order found!",
			"Data":    gin.H{},
		})
		fmt.Println("No such order found!")
		return
	}

	if err := initializer.DB.Preload("Order").Preload("Product").Find(&order, "order_id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't find order items from database!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't find order items from database!")
		return
	}

	if verify.Status == "failed" {
		payment.PaymentStatus = "failed"
		initializer.DB.Save(&payment)

		for _, val := range order {
			val.PaymentStatus = "failed"
			initializer.DB.Save(&val)
		}

		c.JSON(402, gin.H{
			"Status":  "Error!",
			"Code":    402,
			"Message": "Payment failed!",
			"Data":    gin.H{},
		})
		fmt.Println("Payment failed!")
		return
	}

	eror := RazorPaymentVerification(verify.Signature, verify.Order, verify.Payment)
	if eror != nil {
		c.JSON(402, gin.H{
			"Status":  "Error!",
			"Code":    402,
			"Error":   eror.Error(),
			"Message": "Payment verification failed!",
			"Data":    gin.H{},
		})
		fmt.Println("Payment verification failed!")
		return
	}

	payment.TransactionId = verify.Payment
	payment.PaymentStatus = "success"
	erorr := initializer.DB.Save(&payment).Error
	if erorr != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   erorr.Error(),
			"Message": "Couldn't update payment success in database!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't update payment success in database!")
		return
	}

	for _, val := range order {
		initializer.DB.First(&productQuantity, val.ProductId)
		productQuantity.Quantity -= val.Quantity
		if err := initializer.DB.Save(&productQuantity).Error; err != nil {
			fmt.Println("Failed to save updated quantity of products in db:", err)
		}
	}

	for _, val := range order {
		val.PaymentStatus = "success"
		if err := initializer.DB.Save(&val).Error; err != nil {
			fmt.Println("Failed to update payment status for item:", err)
		}
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Payment Successful!",
		"Data":    gin.H{},
	})
	fmt.Println("Payment Successful!")
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := os.Getenv("RAZOR_PAY_SECRET")
	data := orderId + "|" + paymentId

	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("PAYMENT FAILED")
	} else {
		return nil
	}
}
