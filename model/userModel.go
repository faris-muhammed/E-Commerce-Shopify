package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Id        uint          `gorm:"primaryKey"`
	Name      string        `json:"name"`
	Email     string        `gorm:"unique" json:"email"`
	Password  string        `json:"password"`
	Mobile    uint          `json:"mobile"`
	Gender    string        `json:"gender"`
	IsBlocked bool          `gorm:"default:false"`
	IsDeleted bool          `gorm:"default:false"`
	Addresses []UserAddress `gorm:"foreignKey:UserId;references:Id"` // Establishing one-to-many relationship
}

type UserAddress struct {
	gorm.Model
	Id       uint      `gorm:"primaryKey"`
	FullName string    `json:"fullname"`
	Mobile   uint      `json:"mobile"`
	Address  string    `json:"address"`
	Street   string    `json:"street"`
	Landmark string    `json:"landmark"`
	Pincode  uint      `json:"pincode"`
	City     string    `json:"city"`
	UserId   uint      `gorm:"not null"`
	User     UserModel `gorm:"foreignKey:UserId;references:Id"` // Reference to UserModel
}

type Wishlist struct {
	gorm.Model
	ID        uint `gorm:"primary_key"`
	UserID    uint `gorm:"not null"`
	User      UserModel
	ProductID uint `gorm:"not null" json:"product_id"`
	Product   ProductDetails
}

type Cart struct {
	gorm.Model
	Id         uint
	UserId     uint
	User       UserModel
	ProductId  uint
	Product    ProductDetails
	CategoryId uint
	Category   Category
	Quantity   uint
	Price      float64
}

type Order struct {
	gorm.Model
	Id                 uint
	UserId             uint
	User               UserModel
	SellerId           uint
	Seller             SellerModel
	AddressId          uint
	Address            UserAddress
	CouponCode         string
	OrderPaymentMethod string
	OrderAmount        float64
	ShippingCharge     float32
	OrderDate          time.Time
}

type OrderItems struct {
	gorm.Model
	Id            uint `gorm:"primary key"`
	OrderId       uint
	Order         Order
	UserID        uint
	User          UserModel
	SellerId      uint
	Seller        SellerModel
	ProductId     uint
	Product       ProductDetails
	Quantity      uint
	SubTotal      float64
	OrderStatus   string
	PaymentStatus string
}
