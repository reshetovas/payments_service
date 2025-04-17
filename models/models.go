package models

import "time"

type Payment struct {
	ID          int
	Amount      float32
	Description string
	CreatedAt   time.Time
	Currency    string
	ShopID      int
	Address     string
	State       string
	Attempts    int
	Items       []Item
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
