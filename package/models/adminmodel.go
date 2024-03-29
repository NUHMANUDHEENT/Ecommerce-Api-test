package models

import (
	"time"

	"gorm.io/gorm"
)

type Admins struct {
	gorm.Model
	Name     string `gorm:"not null" json:"admin_name"`
	Email    string `gorm:"not null unique" json:"admin_email"`
	Password string `gorm:"not null" json:"admin_password"`
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

type Coupon struct {
	gorm.Model
	Code            string    `gorm:"unique" json:"code"`
	Discount        float64   `json:"discount"`
	CouponCondition int       `json:"condition"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidTo         time.Time `json:"valid_to"`
}
type Offer struct {
	Id           uint
	ProductId    int       `json:"productid"`
	SpecialOffer string    `json:"offer"`
	Discount     float64   `json:"discount"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidTo      time.Time `json:"valid_to"`
}
