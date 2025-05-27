package handlers

import (
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/db"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

// GET /ws/:userId - WebSocket connection
func HandleWebSocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(fiber.StatusUpgradeRequired).SendString("WebSocket upgrade required")
	}

	userID := c.Params("userId")

	return websocket.New(func(conn *websocket.Conn) {
		defer func() {
			db.MongoStorage.RemoveConnection(userID, conn)
			conn.Close()
		}()

		log.Printf("WebSocket connection established for user: %s", userID)

		// Register connection with db
		db.MongoStorage.AddConnection(userID, conn)

		// Send current cart immediately on connection
		cart, err := db.MongoStorage.GetCart(userID)
		if err != nil {
			log.Printf("Error getting cart for user %s: %v", userID, err)
			return
		}

		response := models.WSResponse{
			Type: "cart_updated",
			Cart: cart,
		}

		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error sending initial cart: %v", err)
			return
		}

		// Handle incoming messages
		for {
			var msg models.WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}

			switch msg.Type {
			case "update":
				handleUpdateItem(conn, userID, msg.Data)
			case "remove":
				handleRemoveItem(conn, userID, msg.Data)
			case "clear":
				handleClearCart(conn, userID)
			case "get":
				handleGetCart(conn, userID)
			default:
				sendError(conn, "Unknown message type")
			}
		}
	})(c)
}

func handleUpdateItem(conn *websocket.Conn, userID string, data interface{}) {
	itemData, err := json.Marshal(data)
	if err != nil {
		sendError(conn, "Invalid item data")
		return
	}

	var item models.CartItem
	if err := json.Unmarshal(itemData, &item); err != nil {
		sendError(conn, "Invalid item format")
		return
	}

	// Update in MongoDB (this will trigger polling notification)
	if err := db.MongoStorage.UpdateCartItem(userID, item); err != nil {
		log.Printf("Error updating cart item: %v", err)
		sendError(conn, "Failed to update cart")
		return
	}

	// Polling will handle notifying all connected clients
}

func handleRemoveItem(conn *websocket.Conn, userID string, data interface{}) {
	productID, ok := data.(string)
	if !ok {
		sendError(conn, "Invalid product ID")
		return
	}

	// Remove by setting quantity to 0
	item := models.CartItem{
		ProductID: productID,
		Quantity:  0,
	}

	if err := db.MongoStorage.UpdateCartItem(userID, item); err != nil {
		log.Printf("Error removing cart item: %v", err)
		sendError(conn, "Failed to remove item")
		return
	}

	// Polling will handle notifying all connected clients
}

func handleClearCart(conn *websocket.Conn, userID string) {
	if err := db.MongoStorage.ClearCart(userID); err != nil {
		log.Printf("Error clearing cart: %v", err)
		sendError(conn, "Failed to clear cart")
		return
	}

	// Polling will handle notifying all connected clients
}

func handleGetCart(conn *websocket.Conn, userID string) {
	cart, err := db.MongoStorage.GetCart(userID)
	if err != nil {
		sendError(conn, "Failed to get cart")
		return
	}

	response := models.WSResponse{
		Type: "cart_updated",
		Cart: cart,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending cart update: %v", err)
	}
}

func sendError(conn *websocket.Conn, message string) {
	response := models.WSResponse{
		Type:  "error",
		Error: message,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}
