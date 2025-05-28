package handlers

import (
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/models"
	"github.com/evolza-intern/go_fiber_webapp/api-gateway/services"
	"github.com/gofiber/fiber/v2"
)

func HandleGetNotifications(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/notifications" + "?" + string(c.Context().QueryArgs().QueryString()),
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "notifications", request)
}

func HandleSendNotification(c *fiber.Ctx) error {
	request := models.ProxyRequest{
		Method:  c.Method(),
		URL:     "/notifications",
		Body:    c.Body(),
		Headers: c.GetReqHeaders(),
	}
	return services.ProxyToService(c, "notifications", request)
}
