package entity

import "time"

type Order struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	ProductId  int       `json:"product_id"`
	Qty        int       `json:"qty"`
	TotalPrice int       `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type OrderProduct struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type OrderWithProductMapped struct {
	Id         int           `json:"id"`
	UserId     int           `json:"user_id"`
	Product    *OrderProduct `json:"product"`
	Qty        int           `json:"qty"`
	TotalPrice int           `json:"total_price"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type OrderWithProduct struct {
	Id           int       `json:"id"`
	UserId       int       `json:"user_id"`
	ProductId    int       `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice int       `json:"product_price"`
	Qty          int       `json:"qty"`
	TotalPrice   int       `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
