package models

import "time"

type CartItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// WebSocket message types
type WSMessage struct {
	Type string      `json:"type"` // "update", "remove", "clear", "get"
	Data interface{} `json:"data"`
}

type WSResponse struct {
	Type  string `json:"type"` // "cart_updated", "error"
	Cart  *Cart  `json:"cart,omitempty"`
	Error string `json:"error,omitempty"`
}

// Order structures for checkout (matching order-service)
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

type ShippingAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type CheckoutRequest struct {
	ShippingAddress ShippingAddress `json:"shipping_address"`
	PaymentMethod   string          `json:"payment_method"`
}

type OrderRequest struct {
	UserID          string          `json:"user_id"`
	Items           []OrderItem     `json:"items"`
	TotalAmount     float64         `json:"total_amount"`
	PaymentMethod   string          `json:"payment_method"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}
