package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Gender   string `json:"gender"`
	Phone    int    `gorm:"not null" json:"phone"`
	Blocking bool   `json:"blocking"`
}

type OtpMail struct {
	Id        uint
	Email     string `gorm:"unique" json:"email"`
	Otp       string `gorm:"not null" json:"otp"`
	CreatedAt time.Time
	ExpireAt  time.Time `gorm:"type:timestamp;not null"`
}

type Products struct {
	gorm.Model
	Name        string `gorm:"unique" json:"p_name"`
	Price       uint   `json:"p_price"`
	Size        string `json:"p_size"`
	Color       string `json:"p_color"`
	Quantity    uint   `json:"p_quantity"`
	Description string `json:"p_description"`
	ImagePath1  string
	ImagePath2  string
	ImagePath3  string
	Status      bool `json:"p_blocking"`
	CategoryId  int  `json:"category_id"`
	Category    Category
}
type Category struct {
	gorm.Model
	Category_name        string `gorm:"not null" json:"category_name"`
	Category_description string `gorm:"not null" json:"category_description"`
	Blocking             bool   `gorm:"not null" json:"category_blocking"`
}
type Rating struct {
	gorm.Model
	Users     int `json:"rating_user"`
	ProductId int `gorm:"unique" json:"rating_product"`
	Product   Products
	Value     int `json:"rating_value"`
}
type Review struct {
	Id        uint
	UserId    int `json:"review_user"`
	User      Users
	ProductId uint `json:"review_product"`
	Product   Products
	Review    string `json:"review"`
	Time      string
}
type Address struct {
	gorm.Model
	Address string `json:"user_address"`
	City    string `json:"user_city"`
	State   string `json:"user_state"`
	Pincode int    `json:"user_pincode"`
	Country string `json:"user_country"`
	Phone   int    `json:"user_phone"`
	UserId  int    `json:"user_id"`
	User    Users
}
type Cart struct {
	Id        uint
	UserId    int `json:"user_id"`
	User      Users
	ProductId int
	Product   Products
	Quantity  uint
}
type Coupon struct {
	gorm.Model
	Code      string    `gorm:"unique" json:"code"`
	Discount  float64   `json:"discount"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
}
type Order struct {
	gorm.Model
	UserId        int `json:"order_cart"`
	User          Users
	ProductId     int `json:"order_product"`
	Product       Products
	AddressId     int `json:"order_address"`
	Address       Address
	CouponCode    string `json:"order_coupon"`
	OrderPayment  string `json:"order_payment"`
	OrderQuantity uint
	OrderAmount   float64
	OrderStatus   string
	OrderCancelReason string
}
