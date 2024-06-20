package entity

import "time"

type Transaction struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	OrderId   int       `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
