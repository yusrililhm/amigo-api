package dto

import "time"

type ProductPayload struct {
	Name        string `json:"name" valid:"required~Name can't be empty"`
	Description string `json:"description" valid:"required~Description can't be empty"`
	CategoryId  int    `json:"category_id" valid:"required~Category id can't be empty"`
	Price       int    `json:"price" valid:"required~Price can't be empty"`
	Stock       int    `json:"stock" valid:"required~Stock can't be empty"`
}

type ProductData struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryId  int       `json:"category_id"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Sold        int       `json:"sold"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
