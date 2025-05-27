package dao

import (
	"context"
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type OrderUpdateDAO struct {
	collection *mongo.Collection
}

type UpdateResult struct {
	MatchedCount int64
}

func NewOrderUpdateDAO(collection *mongo.Collection) *OrderUpdateDAO {
	return &OrderUpdateDAO{
		collection: collection,
	}
}

type UpdateFieldResult struct {
	MatchedCount  int64
	ModifiedCount int64
}

func (dao *OrderUpdateDAO) UpdateOrder(ctx context.Context, id primitive.ObjectID, order *models.Order) (*UpdateResult, error) {
	filter := bson.M{
		"order_id": id,
	}

	update := bson.M{
		"$set": bson.M{
			"user_id":           order.UserID,
			"items":             order.Items,
			"total_amount":      order.TotalAmount,
			"status":            order.Status,
			"payment_status":    order.PaymentStatus,
			"payment_method":    order.PaymentMethod,
			"shipping_address":  order.ShippingAddress,
			"order_date":        order.OrderDate.Format(time.RFC3339), // Optional: store as ISO string
			"invoice_generated": order.InvoiceGenerated,
		},
	}

	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &UpdateResult{
		MatchedCount: result.MatchedCount,
	}, nil
}

func (dao *OrderUpdateDAO) UpdateInvoiceGeneratedField(ctx context.Context, orderID primitive.ObjectID) (*UpdateFieldResult, error) {
	filter := bson.M{
		"order_id": orderID,
	}

	update := bson.M{
		"$set": bson.M{
			"invoice_generated": true,
		},
	}

	result, err := dao.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &UpdateFieldResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
	}, nil
}
