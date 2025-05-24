package storage

import (
	"sync"

	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/models"
)

var (
	carts map[string]*models.Cart
	mutex *sync.RWMutex
)

func InitCartStorage() {
	carts = make(map[string]*models.Cart)
	mutex = &sync.RWMutex{}
}

func GetCart(userID string) *models.Cart {
	mutex.RLock()
	defer mutex.RUnlock()

	cart, exists := carts[userID]
	if !exists {
		return &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
			Total:  0.0,
		}
	}
	return cart
}

func SaveCart(userID string, cart *models.Cart) {
	mutex.Lock()
	defer mutex.Unlock()

	carts[userID] = cart
}

func ClearCart(userID string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(carts, userID)
}

func UpdateCartItem(userID string, item models.CartItem) {
	mutex.Lock()
	defer mutex.Unlock()

	cart, exists := carts[userID]
	if !exists {
		cart = &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
			Total:  0.0,
		}
		carts[userID] = cart
	}

	// Find and update existing item or add new one
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			if item.Quantity <= 0 {
				// Remove item if quantity is 0 or less
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i] = item
			}
			found = true
			break
		}
	}

	if !found && item.Quantity > 0 {
		cart.Items = append(cart.Items, item)
	}

	// Recalculate total
	cart.Total = 0.0
	for _, cartItem := range cart.Items {
		cart.Total += cartItem.Subtotal
	}
}
