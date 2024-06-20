package transaction_repo

import "time"

type ProductMapped struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserMapped struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
}

type TransactionWithProductsAndUserMapped struct {
	Id         int           `json:"id"`
	Product    ProductMapped `json:"product"`
	User       UserMapped    `json:"user"`
	Qty        int           `json:"qty"`
	TotalPrice int           `json:"total_price"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
