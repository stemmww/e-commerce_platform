package model

type Product struct {
	ID          string  `json:"id"` // changed from int to string
	Name        string  `json:"name"`
	Category    int     `json:"category"` // changed from CategoryID int to Category string
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"` // changed from int to int32
	Description string  `json:"description"`
}
