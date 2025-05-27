package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/models"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCartStorage struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

// Connection manager for WebSocket clients
type ConnectionManager struct {
	clients map[string]map[*websocket.Conn]bool // userID -> connections
	mutex   sync.RWMutex
}

var connManager = &ConnectionManager{
	clients: make(map[string]map[*websocket.Conn]bool),
}

func NewMongoCartStorage(mongoURI, dbName string) (*MongoCartStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection("carts")

	// Create index on user_id for fast lookups - simplified version
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}

	// Try to create unique index, but don't fail if it already exists
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	_, err = collection.Indexes().CreateOne(ctx2, indexModel)
	if err != nil {
		log.Printf("Warning: Could not create index (may already exist): %v", err)
	}

	storage := &MongoCartStorage{
		client:     client,
		db:         db,
		collection: collection,
	}

	// Start polling instead of change streams (simpler and more compatible)
	go storage.startPolling()

	return storage, nil
}

func (s *MongoCartStorage) GetCart(userID string) (*models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)

	if err == mongo.ErrNoDocuments {
		// Return empty cart if not found
		return &models.Cart{
			UserID:    userID,
			Items:     []models.CartItem{},
			Total:     0.0,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (s *MongoCartStorage) UpdateCartItem(userID string, item models.CartItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calculate subtotal
	item.Subtotal = item.Price * float64(item.Quantity)

	// Get current cart
	currentCart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Update items array
	var updatedItems []models.CartItem
	found := false

	for _, existingItem := range currentCart.Items {
		if existingItem.ProductID == item.ProductID {
			if item.Quantity > 0 {
				updatedItems = append(updatedItems, item)
			}
			// If quantity is 0, we skip adding it (remove item)
			found = true
		} else {
			updatedItems = append(updatedItems, existingItem)
		}
	}

	// If item not found and quantity > 0, add it
	if !found && item.Quantity > 0 {
		updatedItems = append(updatedItems, item)
	}

	// Calculate new total
	var newTotal float64
	for _, cartItem := range updatedItems {
		newTotal += cartItem.Subtotal
	}

	// Update the cart document
	update := bson.M{
		"$set": bson.M{
			"items":      updatedItems,
			"total":      newTotal,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"user_id":    userID,
			"created_at": time.Now(),
		},
	}

	// Use upsert to create if doesn't exist
	opts := options.Update().SetUpsert(true)
	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		update,
		opts,
	)

	return err
}

func (s *MongoCartStorage) ClearCart(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

// Polling approach instead of change streams
func (s *MongoCartStorage) startPolling() {
	ticker := time.NewTicker(1 * time.Second) // Poll every second
	defer ticker.Stop()

	lastUpdate := make(map[string]time.Time)

	for range ticker.C {
		connManager.mutex.RLock()
		userIDs := make([]string, 0, len(connManager.clients))
		for userID := range connManager.clients {
			userIDs = append(userIDs, userID)
		}
		connManager.mutex.RUnlock()

		// Check each user with active connections
		for _, userID := range userIDs {
			cart, err := s.GetCart(userID)
			if err != nil {
				continue
			}

			// Check if cart was updated since last notification
			if lastUpdated, exists := lastUpdate[userID]; !exists || cart.UpdatedAt.After(lastUpdated) {
				s.notifyWebSocketClients(userID, cart)
				lastUpdate[userID] = cart.UpdatedAt
			}
		}
	}
}

func (s *MongoCartStorage) notifyWebSocketClients(userID string, cart *models.Cart) {
	connManager.mutex.RLock()
	clients, exists := connManager.clients[userID]
	connManager.mutex.RUnlock()

	if !exists || len(clients) == 0 {
		return
	}

	// Create WebSocket response using the models.WSResponse struct
	response := models.WSResponse{
		Type: "cart_updated",
		Cart: cart,
	}

	// Send to all connected clients for this user
	for conn := range clients {
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			// Remove failed connection
			s.removeConnection(userID, conn)
		}
	}
}

// WebSocket connection management
func (s *MongoCartStorage) AddConnection(userID string, conn *websocket.Conn) {
	connManager.mutex.Lock()
	defer connManager.mutex.Unlock()

	if connManager.clients[userID] == nil {
		connManager.clients[userID] = make(map[*websocket.Conn]bool)
	}
	connManager.clients[userID][conn] = true

	log.Printf("WebSocket connection added for user: %s", userID)
}

func (s *MongoCartStorage) removeConnection(userID string, conn *websocket.Conn) {
	connManager.mutex.Lock()
	defer connManager.mutex.Unlock()

	if clients, exists := connManager.clients[userID]; exists {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(connManager.clients, userID)
		}
	}

	log.Printf("WebSocket connection removed for user: %s", userID)
}

func (s *MongoCartStorage) RemoveConnection(userID string, conn *websocket.Conn) {
	s.removeConnection(userID, conn)
}

func (s *MongoCartStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}
