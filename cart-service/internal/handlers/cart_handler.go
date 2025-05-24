package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/models"
	"github.com/evolza-intern/go_fiber_webapp/cart-service/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const ORDER_SERVICE_URL = "http://localhost:3001/api/orders"

func InitHandlers() {
	log.Println("Cart handlers initialized")
}

// GET /:userId - Get current cart (fallback)
func GetCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	cart := storage.GetCart(userID)
	cart.UpdatedAt = time.Now()

	return c.JSON(cart)
}

// POST /checkout/:userId - Checkout (saves to order-service)
func CheckoutCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var checkoutReq models.CheckoutRequest
	if err := c.BodyParser(&checkoutReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid checkout data",
		})
	}

	// Get current cart
	cart := storage.GetCart(userID)
	if len(cart.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cart is empty",
		})
	}

	// Convert cart items to order items
	orderItems := make([]models.OrderItem, len(cart.Items))
	for i, item := range cart.Items {
		orderItems[i] = models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		}
	}

	// Create order request
	orderReq := models.OrderRequest{
		UserID:          userID,
		Items:           orderItems,
		TotalAmount:     cart.Total,
		PaymentMethod:   checkoutReq.PaymentMethod,
		ShippingAddress: checkoutReq.ShippingAddress,
	}

	// Send order to order-service
	orderJSON, err := json.Marshal(orderReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to prepare order",
		})
	}

	resp, err := http.Post(ORDER_SERVICE_URL, "application/json", bytes.NewBuffer(orderJSON))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}
	defer resp.Body.Close()

	// if resp.StatusCode != 201 {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"error": "Order service failed to create order",
	// 	})
	// }

	if resp.StatusCode != 201 {
		log.Printf("Order service returned status: %d", resp.StatusCode)
		return c.Status(500).JSON(fiber.Map{
			"error": "Order service failed to create order",
		})
	}

	// Parse order service response
	var orderResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse order response",
		})
	}

	// Clear cart after successful checkout
	storage.ClearCart(userID)

	return c.JSON(fiber.Map{
		"message": "Checkout successful",
		"order":   orderResponse,
	})
}

// GET /ws/:userId - WebSocket connection
func HandleWebSocket(c *fiber.Ctx) error {
	// Check if the request is a WebSocket upgrade
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(fiber.StatusUpgradeRequired).SendString("WebSocket upgrade required")
	}

	userID := c.Params("userId")

	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		log.Printf("WebSocket connection established for user: %s", userID)

		// Send current cart immediately on connection
		cart := storage.GetCart(userID)
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

	// Calculate subtotal
	item.Subtotal = item.Price * float64(item.Quantity)

	storage.UpdateCartItem(userID, item)
	sendCartUpdate(conn, userID)
}

func handleRemoveItem(conn *websocket.Conn, userID string, data interface{}) {
	productID, ok := data.(string)
	if !ok {
		sendError(conn, "Invalid product ID")
		return
	}

	// Remove by setting quantity to 0
	storage.UpdateCartItem(userID, models.CartItem{
		ProductID: productID,
		Quantity:  0,
	})
	sendCartUpdate(conn, userID)
}

func handleClearCart(conn *websocket.Conn, userID string) {
	storage.ClearCart(userID)
	sendCartUpdate(conn, userID)
}

func handleGetCart(conn *websocket.Conn, userID string) {
	sendCartUpdate(conn, userID)
}

func sendCartUpdate(conn *websocket.Conn, userID string) {
	cart := storage.GetCart(userID)
	cart.UpdatedAt = time.Now()

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
