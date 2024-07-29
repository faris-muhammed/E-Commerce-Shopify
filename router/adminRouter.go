package router

import (
	"github.com/gin-gonic/gin"

	controller "main.go/controller/admin"
	"main.go/middleware"
)

var RoleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {
	//======== Landing Page ==========
	r.GET("/adminpage", middleware.AuthMiddleware(RoleAdmin), controller.AdminPage)

	//========== Admin Authentication ==========
	r.GET("/login", controller.AdminLogin)
	r.DELETE("/logout", controller.AdminLogout)

	//============= Category management =============
	r.GET("/category/list", middleware.AuthMiddleware(RoleAdmin), controller.ListCategory)
	r.POST("/category", middleware.AuthMiddleware(RoleAdmin), controller.CreateCategory)
	r.PATCH("/category/edit/:id", middleware.AuthMiddleware(RoleAdmin), controller.EditCategory)
	r.DELETE("/category/delete/:id", middleware.AuthMiddleware(RoleAdmin), controller.DeleteCategory)

	//============= Action's on sellers =============
	r.GET("/sellers", middleware.AuthMiddleware(RoleAdmin), controller.GetAllSellers)
	r.PATCH("/sellers/edit/:id", middleware.AuthMiddleware(RoleAdmin), controller.EditSellerDetails)
	r.PATCH("/sellers/block/:id", middleware.AuthMiddleware(RoleAdmin), controller.BlockSeller)
	r.DELETE("/sellers/delete/:id", middleware.AuthMiddleware(RoleAdmin), controller.DeleteSeller)

	//============= Action's on users =============
	r.GET("/users", middleware.AuthMiddleware(RoleAdmin), controller.GetAllUsers)
	r.PATCH("/users/edit/:id", middleware.AuthMiddleware(RoleAdmin), controller.EditUserDetails)
	r.PATCH("/users/block/:id", middleware.AuthMiddleware(RoleAdmin), controller.BlockUser)
	r.PATCH("/users/delete/:id", middleware.AuthMiddleware(RoleAdmin), controller.DeleteUser)

	//============================= Seller's Order =====================================
	r.GET("/seller/order/:id", middleware.AuthMiddleware(RoleAdmin), controller.ListOrderSeller)
	r.GET("/seller/order/item/:id", middleware.AuthMiddleware(RoleAdmin), controller.ListOrderSeller)

	//============================= Users's Order =====================================
	r.GET("/user/order/:id", middleware.AuthMiddleware(RoleAdmin), controller.ListOrdersUser)
	r.GET("/user/order/item/:id", middleware.AuthMiddleware(RoleAdmin), controller.ListOrderItemsUser)
}
