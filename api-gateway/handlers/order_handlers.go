package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/models"
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/services"
	"github.com/gofiber/fiber/v2"
)

func HandleGetOrders(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/orders" + "?" + string(c.Context().QueryArgs().QueryString()),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "orders", request)
}

func HandleGetOrder(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/orders/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "orders", request)
}

func HandleCreateOrder(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/orders",
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "orders", request)
}

func HandleUpdateOrder(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/orders/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "orders", request)
}

func HandleDeleteOrder(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/orders/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "orders", request)
}
