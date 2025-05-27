package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	OrderID          primitive.ObjectID `bson:"order_id,omitempty" json:"order_id"`
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
