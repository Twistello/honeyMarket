package user

import (
	"time"
)

type User struct {
	Id int64 `db: "id" json:"id"`
	Email string `db: "email" json:"email" validate:"required,email"`
	PasswordHash string `db: "password_hash" json:"-"`
	Role string `db: "role" json:"role" validate:"oneof=admin customer"`
	CreatedAt time.Time `db: "created_at" json:"created_at"`

}