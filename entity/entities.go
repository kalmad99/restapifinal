package entity

import "time"

// Category represents Food Menu Category
type Category struct {
	ID          uint
	Name        string `gorm:"type:varchar(255);not null"`
	Description string
	Image       string `gorm:"type:varchar(255)"`
	Products       []Product `gorm:"many2many:product_categories"`
}

// Role repesents application user roles
type Role struct {
	ID   uint
	Name string `gorm:"type:varchar(255)"`
}

// Item represents food menu items
type Product struct{
	ID uint
	Name string `gorm:"type:varchar(255);not null"`
	ItemType []Category `gorm:"many2many:product_categories"`
	//ItemType string
	Quantity int
	Price float64
	//Seller string
	Description string
	Image string `gorm:"type:varchar(255)"`
	Rating float64
	RatersCount float64
}

// Order represents customer order
type AddToCart struct {
	ID       uint
	PlacedAt time.Time
	UserID   uint
	ItemID   uint
	Quantity uint
}

// User represents application user
//type User struct {
//	ID       uint
//	Name string `gorm:"type:varchar(255);not null"`
//	Email    string `gorm:"type:varchar(255);not null; unique"`
//	Phone    string `gorm:"type:varchar(100);not null; unique"`
//	Password string `gorm:"type:varchar(255)"`
//	//Roles    []Role `gorm:"many2many:user_roles"`
//	//ItemsInCart   []AddToCart
//}
type User struct {
	ID       uint `json:"id"`
	Name string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	//Roles    []Role `gorm:"many2many:user_roles"`
	//ItemsInCart   []AddToCart
}