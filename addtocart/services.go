package carts

import "../entity"

// ItemRepository specifies food menu item related database operations
type CartService interface {
	ItemsinCart() ([]entity.Product, []error)
	IteminCart(id uint) (*entity.Product, []error)
	EditCart(product *entity.Product) (*entity.Product, []error)
	DeleteItem(id uint) (*entity.Product, []error)
	AddtoCart(product *entity.Product) (*entity.Product, []error)
}