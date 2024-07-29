package model

import (
	"time"

	"gorm.io/gorm"
)

type SellerModel struct {
	gorm.Model
	Id          uint   `gorm:"primaryKey"`
	CompanyName string `json:"companyname"`
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	Mobile      uint   `json:"mobile"`
	Pincode     uint   `json:"pincode"`
	Place       string `json:"place"`
	Gst         string `json:"gst"`
	IsBlocked   bool   `gorm:"default:false"`
	IsDeleted   bool   `gorm:"default:false"`
}
type ProductDetails struct {
	gorm.Model
	Id          uint    `gorm:"primarykey"`
	ProductName string  `gorm:"not null" json:"productname"`
	Price       float64 `json:"price"`
	Quantity    uint    `json:"quantity"`
	Size        string  `json:"size"`
	Brand       string  `json:"brand"`
	Barcode     string  `json:"barcode"`
	IsDeleted   bool    `gorm:"default:false"`
	SellerId    uint    `json:"sellerid"`
	Seller      SellerModel
	CategoryId  uint `json:"categoryid"`
	Category    Category
}

type Coupon struct {
	Id             uint      `gorm:"primary_key"`
	Code           string    `gorm:"unique;not null" json:"code"`
	Discount       float64   `gorm:"not null" json:"discount"`
	MaxDiscount    float64   `gorm:"not null" json:"max_discount"`
	MinPurchase    float64   `gorm:"not null" json:"min_purchase"`
	CreatedDate    time.Time `json:"created_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	IsActive       bool      `gorm:"default:false"`
	SellerId       uint
	Seller         SellerModel
}

type OfferProduct struct {
	Id           uint
	SellerId     uint
	Seller       SellerModel
	ProductId    uint `json:"productid"`
	Product      ProductDetails
	SpecialOffer string    `json:"offer"`
	Discount     float64   `json:"discount"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
}

type OfferCategory struct {
	Id           uint
	SellerId     uint
	Seller       SellerModel
	CategoryId   uint `json:"categoryid"`
	Category     Category
	SpecialOffer string    `json:"offer"`
	Discount     float64   `json:"discount"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
}
type Wallet struct {
	gorm.Model
	ID      uint `gorm:"primary_key"`
	UserID  uint
	User    UserModel
	Balance float64 `gorm:"default:0" json:"balance"`
}

type PaymentDetails struct {
	gorm.Model
	UserId        uint
	User          UserModel
	PaymentId     string
	OrderId       uint
	Receipt       uint
	PaymentStatus string
	PaymentAmount float64
	TransactionId string
}
