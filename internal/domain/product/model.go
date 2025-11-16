package product

import (
	"time"
)

type Product struct {
	Id int64 `db: "id" json: "id"`
	Name string `db: "name" json: "name"`
	Category string `db: "category" json: "category"`
	Description string `db: "description" json: "description"`
	Price uint64 `db: "description" json: "description"`
	CreatedAt time.Time `db: "created_at" json: "created_at"`
}