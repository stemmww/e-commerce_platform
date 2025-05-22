package model

type Order struct {
	ID     int         `json:"id"`
	UserID int         `json:"user_id"`
	Status string      `json:"status"`
	Items  []OrderItem `json:"items"` // ✅ Match with proto
}

type OrderItem struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
