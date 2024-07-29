package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"main.go/initializer"
	"main.go/model"
)

//=============== CheckOut ===============

func CheckOut(c *gin.Context) {
	couponCode := ""
	var cartItems []model.Cart
	userId := c.GetUint("userid")
	initializer.DB.Preload("Product").Where("user_id=?", userId).Find(&cartItems)
	if len(cartItems) == 0 {
		c.JSON(409, gin.H{
			"status":  "Fail",
			"message": "please add some items to your cart firstly.",
			"code":    409,
		})
		return
	}
	// ============= check if given payment method and address =============

	paymentMethod := c.Request.PostFormValue("payment")
	Address, _ := strconv.ParseUint(c.Request.PostFormValue("address"), 10, 64)
	if paymentMethod == "" || Address == 0 {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Payment Method and Address are required",
			"code":   400,
		})
		return
	}
	if paymentMethod != "ONLINE" && paymentMethod != "COD" && paymentMethod != "WALLET" {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Give Proper Payment Method ",
			"code":   400,
		})
		return
	}

	// ============= stock check and amount calc ===================
	var Amount float64
	var totalAmount float64
	for _, val := range cartItems {
		productDiscount := OfferDiscountProduct(int(val.ProductId))
		categoryDiscount := OfferDiscountCategory(int(val.Product.CategoryId))
		fmt.Println("Categorydiscount:", categoryDiscount)
		Amount = ((float64(val.Product.Price) - productDiscount - categoryDiscount) * float64(val.Quantity))
		fmt.Println("Amount:", Amount)
		if val.Quantity > val.Product.Quantity {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Insufficent stock for product " + val.Product.ProductName,
				"code":   400,
			})
			return
		}
		totalAmount += Amount
	}

	// ================== coupon validation ===============
	couponCode = c.Request.FormValue("coupon")
	var couponCheck model.Coupon
	var userLimitCheck model.Order
	if couponCode != "" {
		if err := initializer.DB.First(&userLimitCheck, "code", couponCode).Error; err == nil {
			c.JSON(409, gin.H{
				"status": "Fail",
				"error":  "Coupon already used",
				"code":   409,
			})
			return
		}
		if err := initializer.DB.Where(" code=? AND  created_date < ? AND  expiration_date > ? ", couponCode, time.Now(), time.Now()).First(&couponCheck).Error; err != nil {
			c.JSON(200, gin.H{
				"error": "Coupon Not valid",
				"err":   err.Error(),
			})
			return
		} else {
			totalAmount -= couponCheck.Discount
		}
	}

	if totalAmount < couponCheck.MinPurchase {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Fail",
			"error":  "The amount is very low to apply coupon",
			"code":   http.StatusBadRequest,
		})
		return
	}

	// // ================== order id creation =======================
	const charset = "123456789"
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err)
	}
	for i, b := range randomBytes {
		randomBytes[i] = charset[b%byte(len(charset))]
	}
	orderIdstring := string(randomBytes)
	orderId, _ := strconv.Atoi(orderIdstring)
	fmt.Println("-----", orderId)

	//================ Start the transaction ===================

	tx := initializer.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ============== Delivery charges ==============

	var ShippingCharge float64
	if totalAmount < 1000 {
		ShippingCharge = 40
		totalAmount += ShippingCharge
	}

	// ============== COD checking =====================
	if paymentMethod == "COD" {
		if totalAmount > 1000 {
			c.JSON(202, gin.H{
				"status":      "Fail",
				"message":     "Greater than 1000 rupees should not accept COD",
				"totalAmount": totalAmount,
				"code":        202,
			})
			return
		}
	}
	// ================ wallet checking ======================
	if paymentMethod == "WALLET" {
		var walletCheck model.Wallet
		if err := initializer.DB.First(&walletCheck, "user_id=?", userId).Error; err != nil {
			c.JSON(404, gin.H{
				"status": "Fail",
				"error":  "failed to fetch wallet ",
				"err":    err.Error(),
				"code":   404,
			})
			return
		} else if walletCheck.Balance < totalAmount {
			c.JSON(202, gin.H{
				"status": "Fail",
				"error":  "insufficient balance in wallet",
				"code":   202,
			})
			return
		} else {
			// Deduct the amount from the wallet
			walletCheck.Balance -= totalAmount
			if err := tx.Save(&walletCheck).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "Failed to update wallet balance",
					"code":   500,
				})
				return
			}
		}
	}
	// ================= if payment method is online ,redirect to payment actions ==================

	if paymentMethod == "ONLINE" {
		orderid, err := PaymentHandler(orderId, totalAmount)
		if err != nil {
			c.JSON(500, gin.H{
				"status": "Fail",
				"err":    err.Error(),
				"error":  "Failed to create orderId",
				"code":   500,
			})
			tx.Rollback()
			return
		} else {
			c.JSON(200, gin.H{
				"status":      "Success",
				"message":     "please complete the payment",
				"totalAmount": totalAmount,
				"orderId":     orderid,
			})
			err := tx.Create(&model.PaymentDetails{
				OrderId:       uint(orderId),
				PaymentId:     orderid,
				UserId:        userId,
				PaymentAmount: totalAmount,
				PaymentStatus: "pending",
			}).Error
			if err != nil {
				c.JSON(401, gin.H{
					"status": "Fail",
					"error":  "failed to store payment data",
					"err":    err.Error(),
					"code":   401,
				})
				tx.Rollback()
			}
		}
	}
	// ================= insert order details into databse ===================
	for _, val := range cartItems {
		order := model.Order{
			Id:                 uint(orderId),
			UserId:             uint(userId),
			SellerId:           val.Product.SellerId,
			OrderPaymentMethod: paymentMethod,
			AddressId:          uint(Address),
			OrderAmount:        totalAmount,
			ShippingCharge:     float32(ShippingCharge),
			OrderDate:          time.Now(),
			CouponCode:         couponCode,
		}
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{
				"status": "Fail",
				"error":  "failed to place order",
				"err":    err.Error(),
				"code":   500,
			})
			return
		}
	}

	// if err:=tx.Updates()
	// ============ insert order items into database ==================
	for _, val := range cartItems {
		OrderItems := model.OrderItems{
			OrderId:       uint(orderId),
			ProductId:     val.ProductId,
			UserID:        userId,
			SellerId:      val.Product.SellerId,
			Quantity:      val.Quantity,
			SubTotal:      val.Product.Price * float64(val.Quantity),
			OrderStatus:   "pending",
			PaymentStatus: "payment-pending",
		}
		if err := tx.Create(&OrderItems).Error; err != nil {
			tx.Rollback()
			c.JSON(501, gin.H{
				"status": "Fail",
				"error":  "failed to store items details",
				"code":   501,
			})
			return
		}
		// ============= if order is COD manage the stock ============
		if paymentMethod != "ONLINE" {
			var productQuantity model.ProductDetails
			tx.First(&productQuantity, val.ProductId)
			if err := tx.Save(val.Product).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{
					"status": "Fail",
					"error":  "Failed to Update Product Stock",
					"code":   500,
				})
				return
			}
		}
	}
	// =============== delete all items from user cart ==============
	if err := tx.Where("user_id =?", userId).Delete(&model.Cart{}); err.Error != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "faild to delete datas in cart.",
			"code":   400,
		})
		return
	}
	//================= commit transaction whether no error ==================
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "failed to commit transaction",
			"code":   500,
		})
		return
	}
	if paymentMethod != "ONLINE" {
		c.JSON(501, gin.H{
			"status":      "Success",
			"Order":       "Order Placed successfully",
			"payment":     "COD",
			"totalAmount": totalAmount,
			"message":     "Order will arrive with in 4 days",
		})
	}
}
