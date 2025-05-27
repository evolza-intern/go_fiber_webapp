package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/evolza-intern/go_fiber_webapp/order-service/internal/models"
)

func GenerateInvoiceContent(order models.Order) string {
	var sb strings.Builder

	sb.WriteString("===============================================\n")
	sb.WriteString("                 INVOICE                       \n")
	sb.WriteString("===============================================\n\n")

	sb.WriteString(fmt.Sprintf("Invoice OrderID: INV-%s\n", order.OrderID.Hex()))
	sb.WriteString(fmt.Sprintf("Order OrderID: %s\n", order.OrderID.Hex()))
	sb.WriteString(fmt.Sprintf("Order Date: %s\n", order.OrderDate.Format("January 2, 2006")))
	sb.WriteString(fmt.Sprintf("User OrderID: %s\n", order.UserID.Hex()))
	sb.WriteString(fmt.Sprintf("Status: %s\n", strings.Title(order.Status)))
	sb.WriteString(fmt.Sprintf("Payment Status: %s\n", strings.Title(order.PaymentStatus)))
	sb.WriteString(fmt.Sprintf("Payment Method: %s\n\n", strings.Title(order.PaymentMethod)))

	sb.WriteString("SHIPPING ADDRESS:\n")
	sb.WriteString(fmt.Sprintf("%s\n", order.ShippingAddress.Street))
	sb.WriteString(fmt.Sprintf("%s, %s %s\n", order.ShippingAddress.City, order.ShippingAddress.State, order.ShippingAddress.ZipCode))
	sb.WriteString(fmt.Sprintf("%s\n\n", order.ShippingAddress.Country))

	sb.WriteString("ORDER ITEMS:\n")
	sb.WriteString("-----------------------------------------------\n")
	sb.WriteString("Product Name          Qty    Price    Subtotal\n")
	sb.WriteString("-----------------------------------------------\n")

	for _, item := range order.Items {
		sb.WriteString(fmt.Sprintf("%-20s %3d   $%7.2f   $%7.2f\n",
			truncateString(item.ProductName, 20),
			item.Quantity,
			item.Price,
			item.Subtotal))
	}

	sb.WriteString("-----------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("                    TOTAL: $%.2f\n", order.TotalAmount))
	sb.WriteString("===============================================\n\n")

	sb.WriteString(fmt.Sprintf("Generated on: %s\n", time.Now().Format("January 2, 2006 at 3:04 PM")))
	sb.WriteString("Thank you for your business!\n")

	return sb.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
