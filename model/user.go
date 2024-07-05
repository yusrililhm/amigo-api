package model

import (
	"time"
)

type UserData struct {
	Id        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
