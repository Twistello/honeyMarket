package cart

import "time"

type Cart struct {
	Id        uint64    `db:"id" json:"id"`
	UserId    uint64    `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Items     []CartItem `json:"items"`
}

type CartItem struct {
	Id        uint64    `db:"id" json:"id"`
	CartId    uint64    `db:"cart_id" json:"cart_id"`
	ProductId uint64    `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity" validate:"gte=1"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
