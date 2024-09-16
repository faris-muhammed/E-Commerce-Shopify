package router

import (
	"github.com/gin-gonic/gin"
	controller "main.go/controller/user"
	"main.go/middleware"
)

var RoleUser = "User"

func UserGroup(r *gin.RouterGroup) {
	//============= User Authentication =============
	r.POST("/signup", controller.UserSignUp)
	r.POST("/signup/otp", controller.VerifyOTPUser)
	r.POST("/resend/otp", controller.ResendOTP)
	r.GET("/login", controller.UserLogin)
	r.DELETE("/logout", controller.UserLogout)

	//============================ Filter ==============================
	r.GET("/filter", controller.SearchProduct)

	//============================Products & Cart =================================
	r.GET("/listproduct", middleware.AuthMiddleware(RoleUser), controller.ListProducts)
	r.GET("/cart", middleware.AuthMiddleware(RoleUser), controller.ListCart)
	r.POST("/cart/add", middleware.AuthMiddleware(RoleUser), controller.AddCart)
	r.PATCH("/cart/edit/:productId", middleware.AuthMiddleware(RoleUser), controller.EditCart)
	r.DELETE("/cart/delete", middleware.AuthMiddleware(RoleUser), controller.RemoveCart)
	r.POST("/checkout", middleware.AuthMiddleware(RoleUser), controller.CheckOut)

	//============================ Address =============================
	r.GET("/address/list", middleware.AuthMiddleware(RoleUser), controller.ListAddress)
	r.POST("/address", middleware.AuthMiddleware(RoleUser), controller.AddAddress)
	r.PATCH("/address/edit/:id", middleware.AuthMiddleware(RoleUser), controller.EditAddress)
	r.DELETE("/address/delete/:id", middleware.AuthMiddleware(RoleUser), controller.DeleteAddress)

	//============================= Orders =====================================
	r.GET("/order", middleware.AuthMiddleware(RoleUser), controller.ListOrders)
	r.GET("/orderitem/:id", middleware.AuthMiddleware(RoleUser), controller.ListOrderItems)
	r.PATCH("/order/cancel/:id", middleware.AuthMiddleware(RoleUser), controller.CancelOrderItem)

	//========================== Wishlist ===================================
	r.POST("/wishlist/add", middleware.AuthMiddleware(RoleUser), controller.AddToWishlist)
	r.GET("/wishlist", middleware.AuthMiddleware(RoleUser), controller.GetWishlistItems)
	r.PATCH("/wishlist/:id", middleware.AuthMiddleware(RoleUser), controller.RemoveProductFromWishlist)
	r.DELETE("/wishlist/delete", middleware.AuthMiddleware(RoleUser), controller.RemoveWishlist)

	//========================= Wallet =========================
	r.GET("/balance", middleware.AuthMiddleware(RoleUser), controller.WalletBalance)

	//=========================== payment ==========================
	r.GET("/payment", controller.RazorPay)
	r.POST("/payment/confirm", controller.RazorPayVerify)

	//=============== Invoice =================
	r.GET("/order/invoice/:id", middleware.AuthMiddleware(RoleUser), controller.CreateInvoice)
	r.GET("/order/invoice/:id", middleware.AuthMiddleware(RoleUser), controller.CreateInvoice)

}
