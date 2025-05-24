package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName string             `bson:"product_name" json:"product_name"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Price       float64            `bson:"price" json:"price"`
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}

type ShippingAddress struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
	Country string `bson:"country" json:"country"`
}

type Order struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items            []OrderItem        `bson:"items" json:"items"`
	TotalAmount      float64            `bson:"total_amount" json:"total_amount"`
	Status           string             `bson:"status" json:"status"`                 // pending, confirmed, shipped, delivered, cancelled
	PaymentStatus    string             `bson:"payment_status" json:"payment_status"` // pending, paid, failed
	PaymentMethod    string             `bson:"payment_method" json:"payment_method"`
	ShippingAddress  ShippingAddress    `bson:"shipping_address" json:"shipping_address"`
	OrderDate        time.Time          `bson:"order_date" json:"order_date"`
	InvoiceGenerated bool               `bson:"invoice_generated" json:"invoice_generated"`
}
