package models

import "github.com/google/uuid"

type OrderProduct struct {
	ID       int `json:"id,omitempty"`
	Stock    int
	SellerID uuid.UUID
}
