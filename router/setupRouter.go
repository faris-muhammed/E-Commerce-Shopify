package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"https://buynowbazaar.online"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	SetupMiddleware(r)
	SetupRoutes(r)
	SetupTemplates(r)
	return r
}

func SetupMiddleware(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mySession", store))
}

func SetupRoutes(r *gin.Engine) {
	// Grouping
	admin := r.Group("/admin")
	AdminGroup(admin)

	user := r.Group("/user")
	UserGroup(user)

	seller := r.Group("/seller")
	SellerGroup(seller)
}

func SetupTemplates(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
}
