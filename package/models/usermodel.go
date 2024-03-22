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

type Order struct {
	Id                 uint
	UserId             int `json:"order_cart"`
	User               Users
	AddressId          int `json:"order_address"`
	Address            Address
	CouponCode         string `json:"order_coupon"`
	OrderPaymentMethod string `json:"order_payment"`
	OrderAmount        float64
	OrderDate          time.Time
	OrderUpdate        time.Time
}
type OrderItems struct {
	Id                uint `gorm:"primary key"`
	OrderId           uint
	Order             Order
	ProductId         int
	Product           Products
	Quantity          uint
	SubTotal          uint
	OrderStatus       string
	OrderCancelReason string
}
type PaymentDetails struct {
	gorm.Model
	PaymentId     string
	Order_Id      string
	Receipt       uint
	PaymentStatus string
	PaymentAmount int
}
type Wallet struct{
	gorm.Model
	User_id int
	User Users
	Balance float64
}
