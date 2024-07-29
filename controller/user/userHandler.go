package controller

import (
	"fmt"
	"net/http"
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

var RoleUser = "User"

type userDetailSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mobile   uint   `json:"mobile"`
	Gender   string `json:"gender"`
}

// =============== SIGNUP ===============

func UserSignUp(c *gin.Context) {
	var otp string
	var User model.UserModel
	var otpStore model.OTPDetails
	var userDetailsBind userDetailSignUp

	if err := c.Bind(&userDetailsBind); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status": "Fail",
			"error":  "json binding error",
			"code":   http.StatusNotAcceptable,
		})
		return
	}

	if err := initializer.DB.First(&User, "email=?", userDetailsBind.Email).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": "Fail",
			"error":  "Email address already exist",
			"code":   http.StatusConflict,
		})
		return
	}

	otp = controller.GenerateOTP()
	fmt.Println("otp is ----------------", otp, "-----------------")

	if err := controller.SendOTP(userDetailsBind.Email, otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Fail",
			"err":    err.Error(),
			"error":  "failed to send otp",
			"code":   http.StatusBadRequest,
		})
		return
	}

	if result := initializer.DB.First(&otpStore, "email=?", userDetailsBind.Email); result.Error != nil {
		otpStore = model.OTPDetails{
			OTP:       otp,
			Email:     userDetailsBind.Email,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(180 * time.Second),
		}

		if err := initializer.DB.Create(&otpStore); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Fail",
				"error":  "failed to save otp details",
				"code":   http.StatusBadRequest,
			})
			return
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", userDetailsBind.Email).Updates(model.OTPDetails{
			OTP:       otp,
			ExpiresAt: time.Now().Add(180 * time.Second),
		})
		if err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Fail",
				"error":  "Failed to update OTP Details",
				"code":   http.StatusBadRequest,
			})
			return
		}
	}
	userDetails := map[string]interface{}{
		"name":     userDetailsBind.Name,
		"email":    userDetailsBind.Email,
		"password": userDetailsBind.Password,
		"mobile":   userDetailsBind.Mobile,
		"gender":   userDetailsBind.Gender,
	}
	fmt.Println(userDetails)
	session := sessions.Default(c)
	session.Set("signup"+userDetailsBind.Email, userDetails)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Fail",
			"error":  "Failed to save session",
			"err":    err.Error(),
			"code":   http.StatusInternalServerError,
		})
		return
	}

	c.SetCookie("sessionId", "signup"+userDetailsBind.Email, 600, "/", "", false, true)
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "Success",
		"message": "OTP has been sent successfully.",
		"otp":     otp,
		"code":    http.StatusAccepted,
	})
}

// ========== Verify ==========
func VerifyOTPUser(c *gin.Context) {
	var userDataStore model.UserModel
	otp := c.Request.FormValue("otp")
	var existingOTP model.OTPDetails
	if err := initializer.DB.Where("otp = ? AND expires_at > ?", otp, time.Now()).First(&existingOTP).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Fail",
			"error":  "Invalid or expired OTP",
			"code":   http.StatusUnauthorized,
		})
		return
	}
	cookie, err := c.Cookie("sessionId")
	if err != nil || cookie == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "Forbidden",
			"Error":  "Unauthorized Access!",
			"code":   http.StatusForbidden,
		})
		return
	}
	session := sessions.Default(c)
	user := session.Get(cookie)
	fmt.Println("cookie", cookie)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Fail",
			"error":  "User data not found in session",
			"code":   http.StatusNotFound,
		})
		return
	}
	userMap := make(map[string]interface{})
	err = mapstructure.Decode(user, &userMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Fail",
			"error":  "Failed to assert user data to map[string]interface{}",
			"code":   http.StatusInternalServerError,
		})
		return
	}
	mobileStr := fmt.Sprintf("%v", userMap["mobile"])
	gender := fmt.Sprintf("%v", userMap["gender"])
	mobile, _ := strconv.Atoi(mobileStr)
	HashPass, err := bcrypt.GenerateFromPassword([]byte(userMap["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"status": "Fail",
			"error":  "hashing error",
			"code":   http.StatusNotImplemented,
		})
		return
	}
	userDataStore = model.UserModel{
		Name:      userMap["name"].(string),
		Email:     userMap["email"].(string),
		Password:  string(HashPass),
		Mobile:    uint(mobile),
		Gender:    string(gender),
		IsBlocked: false,
	}
	erro := initializer.DB.Create(&userDataStore)
	if erro.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Fail",
			"error":  erro.Error.Error(),
			"code":   http.StatusBadRequest,
		})
		return
	}
	if err := initializer.DB.Delete(&existingOTP).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Fail",
			"error":  "delete data failed",
			"code":   http.StatusBadRequest,
		})
		return
	}
	var userFetchData model.UserModel
	if err := initializer.DB.First(&userFetchData, "email=?", userDataStore.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Fail",
			"error":  "failed to fetch user details for wallet",
			"code":   http.StatusBadRequest,
		})
		return
	}
	fmt.Println("Cookie", cookie)
	session.Delete(cookie)
	session.Save()
	c.SetCookie("sessionId", "", -1, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"status":  "Success",
		"message": "user created successfully",
		"code":    http.StatusCreated,
	})
}

// ============= Resend OTP ================
func ResendOTP(c *gin.Context) {
	var otp string
	var otpStore model.OTPDetails

	cookie, err := c.Cookie("sessionId")
	if err != nil || cookie == "" {
		c.JSON(403, gin.H{
			"status": "Forbidden",
			"Error":  "Unauthorized Access!",
			"code":   http.StatusForbidden,
		})
		return
	}
	session := sessions.Default(c)
	user := session.Get(cookie)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Fail",
			"error":  "User data not found in session",
			"code":   http.StatusNotFound,
		})
		return
	}
	userMap := make(map[string]interface{})
	err = mapstructure.Decode(user, &userMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Fail",
			"error":  "Failed to assert user data to map[string]interface{}",
			"code":   http.StatusInternalServerError,
		})
		return
	}
	otp = controller.GenerateOTP()
	err = controller.SendOTP(userMap["email"].(string), otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
			"code":   http.StatusBadRequest,
		})
		return
	}
	result := initializer.DB.First(&otpStore, "email=?", userMap["email"].(string))
	if result.Error != nil {
		otpStore = model.OTPDetails{
			OTP:       otp,
			Email:     userMap["email"].(string),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(180 * time.Second),
		}
		err := initializer.DB.Create(&otpStore)
		if err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"error":  "failed to store otp",
				"code":   http.StatusBadRequest,
			})
		}
	} else {
		err := initializer.DB.Model(&otpStore).Where("email=?", userMap["email"].(string)).Updates(model.OTPDetails{
			OTP:       otp,
			ExpiresAt: time.Now().Add(180 * time.Second),
		})
		if err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"error":  "failed to update otp",
				"code":   http.StatusBadRequest,
			})
		}
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "OTP has been sent on your registered email id.",
		"otp":     otp,
	})
}

// =============== LOGIN ===============

type UserDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserLogin(c *gin.Context) {
	var UserCheck UserDetails
	var userStore model.UserModel
	if err := c.ShouldBindJSON(&UserCheck); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Binding the data",
			"code":   400,
		})
		return
	}
	if UserCheck.Email == "" || UserCheck.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Email and password are required",
			"code":   401,
		})
		return
	}
	if err := initializer.DB.First(&userStore, "email=?", UserCheck.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "User not found",
			"code":   401,
		})
		return
	}
	// Check if the user account is blocked
	if userStore.IsBlocked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Your account is blocked. Please contact support.",
			"code":   403,
		})
		return
	}
	// Compare the password
	err := bcrypt.CompareHashAndPassword([]byte(userStore.Password), []byte(UserCheck.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Invalid Username or password",
			"code":   500,
		})
		return
	}
	token, err := middleware.GenerateToken(userStore.Id, userStore.Email, RoleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"error":  "Failed to generate JWT token",
			"code":   500,
		})
		return
	}
	fmt.Println(userStore.Id, userStore.Email)
	c.SetCookie("jwtTokenUser", token, int(time.Now().Add(1*time.Hour).Unix()), "/", "", false, true)
	fmt.Println("Token", token)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"code":   200,
	})
}

//=============== LOGOUT ===============

func UserLogout(c *gin.Context) {
	// Clear JWT token cookie
	c.SetCookie("jwtTokenUser", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged out successfully",
		"code":    http.StatusOK,
	})
}
