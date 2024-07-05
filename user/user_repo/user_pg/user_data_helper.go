package user_pg

import (
	"database/sql"
	"time"
)

type userData struct {
	Id        int            `json:"id"`
	FullName  string         `json:"full_name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      string         `json:"role"`
	Address   sql.NullString `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
