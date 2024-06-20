package entity

import "time"

type Product struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryId  int       `json:"category_id"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Sold        int       `json:"sold"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
