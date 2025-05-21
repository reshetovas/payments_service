package models

import (
	"time"
)

type Payment struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float32   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Currency    string    `json:"currency"`
	ShopID      int       `json:"shop_id"`
	Address     string    `json:"address"`
	State       string    `json:"state"`
	Attempts    int       `json:"attempts"`
	Items       []Item    `json:"items"`
}

type Item struct {
	ID        int
	PaymentID int
	Name      string
	Price     float32
	Quantity  int
}

type Bonus struct {
	ID        int
	PaymentID int
	Amount    float32
}

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}
