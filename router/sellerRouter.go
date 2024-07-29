package router

import (
	"github.com/gin-gonic/gin"
	controller "main.go/controller/seller"

	"main.go/middleware"
)

var RoleSeller = "Seller"

func SellerGroup(r *gin.RouterGroup) {
	//============== Seller Authentication ==============
	r.POST("/signup", controller.SellerSignUp)
	r.POST("/signup/otp", controller.VerifyOTPSeller)
	r.POST("/resend/otp", controller.ResendOTP)
	r.GET("/login", controller.SellerLogin)
	r.DELETE("/logout", controller.SellerLogout)

	//============== Product ==============
	r.GET("/product/list", middleware.AuthMiddleware(RoleSeller), controller.ListProduct)
	r.POST("/product", middleware.AuthMiddleware(RoleSeller), controller.AddProduct)
	r.PATCH("/product/edit/:id", middleware.AuthMiddleware(RoleSeller), controller.EditProduct)
	r.PATCH("/product/delete/:id", middleware.AuthMiddleware(RoleSeller), controller.SoftDeleteProduct)
	r.PATCH("/product/recover/:id", middleware.AuthMiddleware(RoleSeller), controller.RecoverDeleteProduct)

	//============================= Orders =====================================
	r.GET("/order", middleware.AuthMiddleware(RoleSeller), controller.ListOrders)
	r.PATCH("/order/deliver/:id", middleware.AuthMiddleware(RoleSeller), controller.DeliverOrder)
	r.PATCH("/order/cancel/:id", middleware.AuthMiddleware(RoleSeller), controller.CancelOrder)

	//============================ Coupon ====================================
	r.GET("/coupon", middleware.AuthMiddleware(RoleSeller), controller.CouponView)
	r.POST("/coupon/create", middleware.AuthMiddleware(RoleSeller), controller.CouponCreate)
	r.POST("/coupon/delete/:id", middleware.AuthMiddleware(RoleSeller), controller.CouponDelete)

	// =================== offer management =====================
	//======  Product ==========
	r.GET("/offer", middleware.AuthMiddleware(RoleSeller), controller.OfferProductList)
	r.POST("/offer/add", middleware.AuthMiddleware(RoleSeller), controller.OfferProductAdd)
	r.DELETE("/offer/delete/:id", middleware.AuthMiddleware(RoleSeller), controller.OfferProductDelete)

	//===== Category ===========
	r.GET("/listcategory", middleware.AuthMiddleware(RoleSeller), controller.ListCategory)
	r.GET("/offercategory", middleware.AuthMiddleware(RoleSeller), controller.OfferCategoryList)
	r.POST("/offercategory/add", middleware.AuthMiddleware(RoleSeller), controller.OfferCategoryAdd)
	r.DELETE("/offercategory/delete/:id", middleware.AuthMiddleware(RoleSeller), controller.OfferCategoryDelete)

	// ===================== sales report =========================
	r.GET("/sales/report", middleware.AuthMiddleware(RoleSeller), controller.SalesReport)
	r.GET("/sales/report/excel", middleware.AuthMiddleware(RoleSeller), controller.SalesReportExcel)
	r.GET("/sales/report/pdf", middleware.AuthMiddleware(RoleSeller), controller.SalesReportPDF)

	// ===================== Best selling ========================
	r.GET("/bestselling", middleware.AuthMiddleware(RoleSeller), controller.BestSelling)

}
