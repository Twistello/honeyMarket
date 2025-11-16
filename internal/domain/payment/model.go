package payment

import (
	"time"
)

type Payment struct {
	Id int64 `db: "id" json: "id"`
	Amount uint64 `db: "amount" json: "amount"`
	OrderId int64 `db: "order_id" json: "order_id"`
	Status string `db: "status" json: "status"`
	TransactionId string `db: "transaction_id" json: "transaction_id"`
	Provider string `db: "provider" json: "provider"`
	CreatedAt time.Time `db: "created_at" json: "created_at"`
}