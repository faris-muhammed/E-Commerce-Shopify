package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email  string `json:"username"`
	Role   string `json:"role"`
	UserId uint
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRETKEY"))

func JwtToken(c *gin.Context, userId uint, email string, role string) {
	tokenString, err := GenerateToken(userId, email, role)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "Fail",
			"error":  "Failed to generate JWT token",
			"code":   400,
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "Success",
		"code":   200,
		"token":  tokenString,
	})
}

func GenerateToken(userId uint, email string, role string) (string, error) {
	claims := Claims{
		Email:  email,
		Role:   role,
		UserId: uint(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwtToken" + requiredRole)
		fmt.Println("TokenString", tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Can't find cookie",
				"code":    401,
			})
			c.Abort()
			return
		}
		if tokenString == "" {
			c.JSON(400, gin.H{
				"status":  "Bad Request",
				"message": "Empty token string.",
				"code":    400,
			})
			c.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Tokenclaims", token.Claims)
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			fmt.Println("cookie error:", err)
			c.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Invalid or expired JWT Token.",
				"code":    401,
			})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			fmt.Println("req", requiredRole, claims.Role)
			c.JSON(403, gin.H{
				"status": "Forbidden",
				"error":  "Insufficient permissions",
				"code":   403,
			})
			c.Abort()
			return
		}
		c.Set("userid", claims.UserId)
		c.Next()
	}
}
