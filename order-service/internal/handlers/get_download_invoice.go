package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

// DownloadInvoice handles the download of invoice files for orders by validating the filename and serving the file
func DownloadInvoice(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Security: Prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid filename",
		})
	}

	// Validate filename format (should be invoice_<orderid>.txt)
	if !strings.HasPrefix(filename, "invoice_") || !strings.HasSuffix(filename, ".txt") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid invoice filename format",
		})
	}

	filepath := fmt.Sprintf("./invoices/%s", filename)

	// Check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "Invoice file not found",
		})
	}

	// Set headers for file download
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Set("Content-Type", "text/plain")

	return c.SendFile(filepath)
}
