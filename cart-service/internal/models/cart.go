package models

import "time"

type CartItem struct {
	ProductID   string  `bson:"product_id" json:"product_id"`
	ProductName string  `bson:"product_name" json:"product_name"`
	Price       float64 `bson:"price" json:"price"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Subtotal    float64 `bson:"subtotal" json:"subtotal"`
}

type Cart struct {
	ID        interface{} `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string      `bson:"user_id" json:"user_id"`
	Items     []CartItem  `bson:"items" json:"items"`
	Total     float64     `bson:"total" json:"total"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
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
