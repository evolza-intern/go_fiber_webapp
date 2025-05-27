package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/products-service/internal/dao"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var productCollection *mongo.Collection

func InitProductHandler(collection *mongo.Collection) {
	productCollection = collection
}

type ProductHandlers struct {
	createDAO *dao.ProductCreateDAO
	readDAO   *dao.ProductReadDAO
	updateDAO *dao.ProductUpdateDAO
	deleteDAO *dao.ProductDeleteDAO
}

func NewProductHandlers(collection *mongo.Collection) *ProductHandlers {
	return &ProductHandlers{
		createDAO: dao.NewProductCreateDAO(collection),
		readDAO:   dao.NewProductReadDAO(collection),
		updateDAO: dao.NewProductUpdateDAO(collection),
		deleteDAO: dao.NewProductDeleteDAO(collection),
	}
}
