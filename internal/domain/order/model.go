package order

import (
	"time"
)

type Order struct {
	Id int64 `db:"id" json: "id"`
	UserId int64  `db:"user_id" json: "user_id"`
	Total uint64 `db:"total" json: "total"`
	Status string `db:"status" json: "status"`
	CreatedAt time.Time `db:"created_at" json: "created_at"`
}

type OrderItem struct {
	Id int64 `db:"id" json: "id"`
	OrderId int64 `db:"order_id" json: "order_id"`
	ProductID int64 `db:"product_id" json: "product_id"`
	Quantity int `db:"quantity" json: "quantity"`
	Price uint64 `db:"price" json: "price"`
}