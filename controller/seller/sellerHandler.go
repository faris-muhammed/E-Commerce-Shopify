package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"main.go/controller"
	"main.go/initializer"
	"main.go/middleware"
	"main.go/model"
)

var RoleSeller = "Seller"

// =============== SIGNUP ===============
type sellerDetailSignUp struct {
	CompanyName string `json:"companyName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Mobile      uint   `json:"mobile"`
	Pincode     uint   `json:"pincode"`
	Place       string `json:"place"`
	Gst         string `json:"gst"`
}

func SellerSignUp(c *gin.Context) {
	var otp string
	var Seller model.SellerModel
	var otpStore model.OTPDetails
	var sellerDetailsBind sellerDetailSignUp

	if err := c.Bind(&sellerDetailsBind); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "json binding error",
			"code":   400,
		})
		return
	}

	if err := initializer.DB.First(&Seller, "email=?", sellerDetailsBind.Email).Error; err == nil {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "Email address already exist",
			"code":   409,
		})
		return
	}

	otp = controller.GenerateOTP()
	fmt.Println("otp is ----------------", otp, "-----------------")

	if err := controller.SendOTP(sellerDetailsBind.Email, otp); err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "failed to send otp",
			"code":   400,
		})
		return
	}

	if result := initializer.DB.First(&otpStore, "email=?", sellerDetailsBind.Email); result.Error != nil {
		otpStore = model.OTPDetails{
			OTP:       otp,
			Email:     sellerDetailsBind.Email,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(180 * time.Second),
		}

		if err := initializer.DB.Create(&otpStore); err.Error != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Failed to save otp details",
				"code":   400,
			})
			return
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", sellerDetailsBind.Email).Updates(model.OTPDetails{
			OTP:       otp,
			ExpiresAt: time.Now().Add(180 * time.Second),
		})
		if err.Error != nil {
			c.JSON(400, gin.H{
				"status": "Fail",
				"error":  "Failed to update OTP Details",
				"code":   400,
			})
			return
		}
	}
	sellerDetails := map[string]interface{}{
		"companyName": sellerDetailsBind.CompanyName,
		"email":       sellerDetailsBind.Email,
		"password":    sellerDetailsBind.Password,
		"mobile":      sellerDetailsBind.Mobile,
		"pincode":     sellerDetailsBind.Pincode,
		"place":       sellerDetailsBind.Place,
		"gst":         sellerDetailsBind.Gst,
	}
	fmt.Println(sellerDetails)
	session := sessions.Default(c)
	session.Set("signup"+sellerDetailsBind.Email, sellerDetails)
	if err := session.Save(); err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to save session",
			"code":   500,
		})
		return
	}

	c.SetCookie("sessionId", "signup"+sellerDetailsBind.Email, 600, "/", "", false, true)
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "OTP has been sent successfully.",
		"otp":     otp,
		"code":    200,
	})
}

// ========== Verify ==========
func VerifyOTPSeller(c *gin.Context) {
	var sellerDataStore model.SellerModel
	otp := c.Request.FormValue("otp")
	var existingOTP model.OTPDetails
	if err := initializer.DB.Where("otp = ? AND expires_at > ?", otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Invalid or expired OTP",
			"code":   401,
		})
		return
	}
	cookie, err := c.Cookie("sessionId")
	if err != nil || cookie == "" {
		c.JSON(401, gin.H{
			"status": "Forbidden",
			"Error":  "Unauthorized Access!",
			"code":   401,
		})
		return
	}
	session := sessions.Default(c)
	seller := session.Get(cookie)
	if seller == nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "User data not found in session",
			"code":   404,
		})
		return
	}
	sellerMap := make(map[string]interface{})
	err = mapstructure.Decode(seller, &sellerMap)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to assert seller data to map[string]interface{}",
			"code":   500,
		})
		return
	}
	mobileStr := fmt.Sprintf("%v", sellerMap["mobile"])
	pincodeStr := fmt.Sprintf("%v", sellerMap["pincode"])
	mobile, _ := strconv.Atoi(mobileStr)
	pincode, _ := strconv.Atoi(pincodeStr)
	HashPass, err := bcrypt.GenerateFromPassword([]byte(sellerMap["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "hashing error",
			"code":   409,
		})
		return
	}
	sellerDataStore = model.SellerModel{
		CompanyName: sellerMap["companyName"].(string),
		Email:       sellerMap["email"].(string),
		Password:    string(HashPass),
		Mobile:      uint(mobile),
		Pincode:     uint(pincode),
		Place:       sellerMap["place"].(string),
		Gst:         sellerMap["gst"].(string),
		IsBlocked:   false,
		IsDeleted:   false,
	}
	erro := initializer.DB.Create(&sellerDataStore)
	if erro.Error != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to create",
			"code":   400,
		})
		return
	}
	if err := initializer.DB.Delete(&existingOTP).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to delete existing OTP",
			"code":   400,
		})
		return
	}
	var sellerFetchData model.SellerModel
	if err := initializer.DB.First(&sellerFetchData, "email=?", sellerDataStore.Email).Error; err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to fetch seller details for wallet",
			"code":   400,
		})
		return
	}
	session.Delete(cookie)
	session.Save()
	c.SetCookie("sessionId", "", -1, "/", "", false, true)
	c.JSON(201, gin.H{
		"status":  "Success",
		"message": "seller created successfully",
		"code":    201,
	})
}

// ============= Resend OTP ================
func ResendOTP(c *gin.Context) {
	var otp string
	var otpStore model.OTPDetails

	cookie, err := c.Cookie("sessionId")
	if err != nil || cookie == "" {
		c.JSON(401, gin.H{
			"status": "Forbidden",
			"Error":  "Unauthorized Access!",
			"code":   401,
		})
		return
	}
	session := sessions.Default(c)
	seller := session.Get(cookie)
	if seller == nil {
		c.JSON(404, gin.H{
			"status": "Fail",
			"error":  "User data not found in session",
			"code":   404,
		})
		return
	}
	sellerMap := make(map[string]interface{})
	err = mapstructure.Decode(seller, &sellerMap)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to assert seller data to map[string]interface{}",
			"code":   500,
		})
		return
	}
	otp = controller.GenerateOTP()
	err = controller.SendOTP(sellerMap["email"].(string), otp)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"error":  err.Error(),
			"code":   400,
		})
		return
	}
	result := initializer.DB.First(&otpStore, "email=?", sellerMap["email"].(string))
	if result.Error != nil {
		otpStore = model.OTPDetails{
			OTP:       otp,
			Email:     sellerMap["email"].(string),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(180 * time.Second),
		}
		err := initializer.DB.Create(&otpStore)
		if err.Error != nil {
			c.JSON(400, gin.H{
				"status": "fail",
				"error":  "failed to store otp",
				"code":   400,
			})
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", sellerMap["email"].(string)).Updates(model.OTPDetails{
			OTP:       otp,
			ExpiresAt: time.Now().Add(180 * time.Second),
		})
		if err.Error != nil {
			c.JSON(400, gin.H{
				"status": "fail",
				"error":  "failed to update otp",
				"code":   400,
			})
		}
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "OTP has been sent on your registered email_id",
		"code":    200,
		"otp":     otp,
	})
}

// =========== LOGIN ===========
type SellerDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SellerLogin(c *gin.Context) {
	var SellerCheck SellerDetails
	var sellerstore model.SellerModel
	// Binding the json data from requested URL to sellerDetails struct
	if err := c.ShouldBindJSON(&SellerCheck); err != nil {
		c.JSON(401, gin.H{
			"status":  "Fail",
			"message": "Binding the data",
			"error":   err.Error(),
			"code":    400,
		})
		return
	}
	// Checking the required credentials given or not
	if SellerCheck.Email == "" || SellerCheck.Password == "" {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Email and password are required",
			"code":   401,
		})
		return
	}
	//
	if err := initializer.DB.First(&sellerstore, "email=?", SellerCheck.Email).Error; err != nil {
		c.JSON(401, gin.H{
			"status":  "Fail",
			"message": "Seller not found",
			"error":   err.Error(),
			"code":    401,
		})
		return
	}
	// Check if the seller account is blocked
	if sellerstore.IsBlocked {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Your account is blocked. Please contact support.",
			"code":   401,
		})
		return
	}
	// Check if the seller account is deleted
	if sellerstore.IsDeleted {
		c.JSON(401, gin.H{
			"status": "Fail",
			"error":  "Your account is deleted",
			"code":   401,
		})
		return
	}
	// Comparing the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(sellerstore.Password), []byte(SellerCheck.Password))
	if err != nil {
		c.JSON(409, gin.H{
			"status": "Fail",
			"error":  "Invalid Username or password",
			"code":   409,
		})
	}
	// Generating token
	token, err := middleware.GenerateToken(sellerstore.Id, sellerstore.Email, RoleSeller)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Fail",
			"error":  "Failed to generate JWT token",
			"code":   500,
		})
		return
	}
	// Setting the token inside cookie
	c.SetCookie("jwtTokenSeller", token, int(time.Now().Add(1*time.Hour).Unix()), "/", "", false, true)

	c.JSON(200, gin.H{
		"status": "Success",
		"token":  token,
		"code":   200,
	})
}

//========== LOGOUT ==========

func SellerLogout(c *gin.Context) {
	// Clear JWT token cookie
	c.SetCookie("jwtTokenSeller", "", -1, "/", "", false, true)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Seller logged out successfully",
		"code":    200,
	})
}
