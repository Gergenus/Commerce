package models

type OrderProduct struct {
	ID          int     `json:"id,omitempty"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	SellerID    string  `json:"seller_id,omitempty"`
	CategoryID  int     `json:"category_id"`
	Quantity    int     `json:"quantity"`
}
