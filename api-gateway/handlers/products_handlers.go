package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/models"
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/services"
	"github.com/gofiber/fiber/v2"
)

func HandleGetProducts(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/products" + "?" + string(c.Context().QueryArgs().QueryString()),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "products", request)
}

func HandleGetProduct(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/products/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "products", request)
}

func HandleCreateProduct(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/products",
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "products", request)
}

func HandleUpdateProduct(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/products/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "products", request)
}

func HandleDeleteProduct(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/products/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "products", request)
}
