package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/models"
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/services"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/users" + "?" + string(c.Context().QueryArgs().QueryString()),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "users", request)
}

func HandleGetUser(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/users/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "users", request)
}

func HandleCreateUser(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/users",
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "users", request)
}

func HandleUpdateUser(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/users/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "users", request)
}

func HandleDeleteUser(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/users/" + c.Params("id"),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "users", request)
}
