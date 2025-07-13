package models

type AddCartRequest struct {
	ProductId string `json:"product_id"`
	Stock     int    `json:"stock"`
}
