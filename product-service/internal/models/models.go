package models

type Product struct {
	ID          int     `json:"id,omitempty"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	SellerID    int     `json:"seller_id"`
	CategoryID  int     `json:"category_id"`
}

type Category struct {
	Category string `json:"category"`
}

type AddStockRequest struct {
	SellerID  int `json:"seller_id"`
	ProductID int `json:"product_id"`
	Number    int `json:"number"`
}
