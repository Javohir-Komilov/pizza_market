package handlers

import (
	"log"
	"net/http"

	"pizza/db"
	"pizza/templates"

	"github.com/labstack/echo/v4"
)

func PlaceOrderHandler(c echo.Context) error {
	user := c.Get("user").(db.User)
	ctx := c.Request().Context()
	cartItems, err := db.Q.GetCartItems(ctx, user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch cart items")
	}

	var total float64

	for _, item := range cartItems {
        menuItem, err := db.Q.GetMenuItemById(ctx, item.MenuItemID)
        if err != nil {
            log.Printf("Error fetching menu item: %v", err)
            return c.String(http.StatusInternalServerError, "Failed to fetch menu items")
        }
        total += menuItem.Price * float64(item.Quantity)	
    }

	// Create order
	order, err := db.Q.CreateOrder(ctx, db.CreateOrderParams{
		UserID:     user.ID,
		TotalPrice: total,
		Status:     "Pending",
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create order")
	}

	for _, item := range cartItems {
		_, err := db.Q.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:    order.ID,
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
		})
		if err != nil {
			
			return c.String(http.StatusInternalServerError, "Failed to add items to order")
		}
	}

	// Clear the cart
	err = db.Q.ClearCart(ctx, user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to clear cart")
	}

	return c.Redirect(http.StatusSeeOther, "/user/orders")
}

func UserOrdersHandler(c echo.Context) error {
	user := c.Get("user").(db.User)
	orders, err := db.Q.ListOrdersByUser(c.Request().Context(), user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch orders")
	}

	return templates.Render(c, templates.Orders(orders))
}

// func CancelOrderHandler(c echo.Context) error {
//     user := c.Get("user").(db.User)
//     orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
//     if err != nil {
//         return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Invalid order ID"})
//     }

//     order, err := db.Q.GetOrder(c.Request().Context(), orderID)
//     if err != nil || order.UserID != user.ID || order.Status != "Pending" {
//         return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Cannot cancel this order"})
//     }

//     err = db.Q.UpdateOrderStatus(c.Request().Context(), db.UpdateOrderStatusParams{
//         ID:     orderID,
//         Status: "Cancelled",
//     })
//     if err != nil {
//         return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "Failed to cancel order"})
//     }

//     return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
// }

// func ReorderHandler(c echo.Context) error {
//     user := c.Get("user").(db.User)
//     orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
//     if err != nil {
//         return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Invalid order ID"})
//     }

//     oldOrder, err := db.Q.GetOrder(c.Request().Context(), orderID)
//     if err != nil || oldOrder.UserID != user.ID {
//         return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Cannot reorder this order"})
//     }

//     newOrder, err := db.Q.CreateOrder(c.Request().Context(), db.CreateOrderParams{
//         UserID:     user.ID,
//         TotalPrice: oldOrder.TotalPrice,
//         Status:     "Pending",
//     })
//     if err != nil {
//         return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "Failed to create new order"})
//     }

//     oldOrderItems, err := db.Q.GetOrderItems(c.Request().Context(), orderID)
//     if err != nil {
//         return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "Failed to get old order items"})
//     }

//     for _, item := range oldOrderItems {
//         _, err := db.Q.CreateOrderItem(c.Request().Context(), db.CreateOrderItemParams{
//             OrderID:    newOrder.ID,
//             MenuItemID: item.MenuItemID,
//             Quantity:   item.Quantity,
//         })
//         if err != nil {
//             return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "Failed to add items to new order"})
//         }
//     }

//     return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
// }