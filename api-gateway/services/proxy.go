package services

import (
	"bytes"
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/models"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var serviceRegistry = map[string]models.Service{
	"users": {
		Name:    "Users Service",
		BaseURL: getEnv("USERS_SERVICE_URL", "http://localhost:8001"),
		Timeout: getEnvInt("USERS_SERVICE_TIMEOUT", 30),
	},
	"products": {
		Name:    "Products Service",
		BaseURL: getEnv("PRODUCTS_SERVICE_URL", "http://localhost:8002"),
		Timeout: getEnvInt("PRODUCTS_SERVICE_TIMEOUT", 30),
	},
	"orders": {
		Name:    "Orders Service",
		BaseURL: getEnv("ORDERS_SERVICE_URL", "http://localhost:8003"),
		Timeout: getEnvInt("ORDERS_SERVICE_TIMEOUT", 30),
	},
	"notifications": {
		Name:    "Notifications Service",
		BaseURL: getEnv("NOTIFICATIONS_SERVICE_URL", "http://localhost:8004"),
		Timeout: getEnvInt("NOTIFICATIONS_SERVICE_TIMEOUT", 30),
	},
}

func ProxyToService(c *fiber.Ctx, serviceName string, request models.ProxyRequest) error {
	service, exists := serviceRegistry[serviceName]
	if !exists {
		return c.Status(404).JSON(models.ErrorResponse{
			Error:   "Service not found",
			Message: "The requested service is not available",
			Code:    404,
		})
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(service.Timeout) * time.Second,
	}

	// Prepare request URL
	url := service.BaseURL + request.URL

	// Create request
	req, err := http.NewRequest(request.Method, url, bytes.NewReader(request.Body))
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error:   "Request creation failed",
			Message: err.Error(),
			Code:    500,
		})
	}

	// Copy headers
	for key, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(503).JSON(models.ErrorResponse{
			Error:   "Service unavailable",
			Message: "Failed to connect to " + service.Name,
			Code:    503,
		})
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error:   "Response reading failed",
			Message: err.Error(),
			Code:    500,
		})
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header.Add(key, value)
		}
	}

	// Return response
	return c.Status(resp.StatusCode).Send(body)
}

func GetServiceRegistry() map[string]models.Service {
	return serviceRegistry
}

func AddService(name string, service models.Service) {
	serviceRegistry[name] = service
}

func RemoveService(name string) {
	delete(serviceRegistry, name)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
