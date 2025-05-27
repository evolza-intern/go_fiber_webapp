package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/dao"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var orderCollection *mongo.Collection

func InitOrderHandler(collection *mongo.Collection) {
	orderCollection = collection
}

type OrderHandlers struct {
	createDAO *dao.OrderCreateDAO
	readDAO   *dao.OrderReadDAO
	updateDAO *dao.OrderUpdateDAO
	deleteDAO *dao.OrderDeleteDAO
}

func NewOrderHandlers(collection *mongo.Collection) *OrderHandlers {
	return &OrderHandlers{
		createDAO: dao.NewOrderCreateDAO(collection),
		readDAO:   dao.NewOrderReadDAO(collection),
		updateDAO: dao.NewOrderUpdateDAO(collection),
		deleteDAO: dao.NewOrderDeleteDAO(collection),
	}
}
