package handlers

import (
	"log"
	"net/http"
	"pizza/db"
	"pizza/models"
	"pizza/templates"
	"strconv"

	"github.com/labstack/echo/v4"
)



func CartHandler(c echo.Context) error {
    user, ok := c.Get("user").(db.User)
    if !ok {
        return c.String(http.StatusInternalServerError, "User not found")
    }
    ctx := c.Request().Context()

    cartItems, err := db.Q.GetCartItems(ctx, user.ID)
    if err != nil {
        log.Printf("Error fetching cart items: %v", err)
        return c.String(http.StatusInternalServerError, "Failed to fetch cart items")
    }

    var totalPrice float64
    var cartItemsInfo []models.CartItemInfo

    for _, item := range cartItems {
        menuItem, err := db.Q.GetMenuItemById(ctx, item.MenuItemID)
        if err != nil {
            log.Printf("Error fetching menu item: %v", err)
            return c.String(http.StatusInternalServerError, "Failed to fetch menu items")
        }
        subtotal := menuItem.Price * float64(item.Quantity)
        cartItemsInfo = append(cartItemsInfo, models.CartItemInfo{
            MenuItem: menuItem,
            Quantity: item.Quantity,
            Subtotal: subtotal,
        })
        totalPrice += subtotal
    }

    categories, err := db.Q.GetCategories(ctx)
    if err != nil {
        log.Printf("Error fetching categories: %v", err)
        return c.String(http.StatusInternalServerError, "Failed to fetch categories")
    }

    cartData := struct {
        CartItems  []models.CartItemInfo
        TotalPrice float64
        Categories []db.Category
        User       db.User
    }{
        CartItems:  cartItemsInfo,
        TotalPrice: totalPrice,
        Categories: categories,
        User:       user,
    }

    return templates.Render(c, templates.Cart(cartData))
}

func AddToCartHandler(c echo.Context) error {
    user, ok := c.Get("user").(db.User)
    if !ok {
        return c.String(http.StatusInternalServerError, "User not found")
    }
    ctx := c.Request().Context()

    itemID, err := strconv.ParseInt(c.FormValue("item_id"), 10, 64)
    if err != nil {
        return c.String(http.StatusBadRequest, "Invalid item ID")
    }

    quantity, err := strconv.ParseInt(c.FormValue("quantity"), 10, 64)
    if err != nil {
        return c.String(http.StatusBadRequest, "Invalid item ID")
    }

    _, err = db.Q.GetMenuItemById(ctx, itemID)
    if err != nil {
        log.Printf("Error fetching menu item: %v", err)
        return c.String(http.StatusNotFound, "Menu item not found")
    }

    cartItems, err := db.Q.GetCartItems(ctx, user.ID)
    if err != nil {
        log.Printf("Error fetching cart items: %v", err)
        return c.String(http.StatusInternalServerError, "Failed to fetch cart items")
    }

    itemExists := false
    var existingItem db.CartItem
    for _, item := range cartItems {
        if item.MenuItemID == itemID {
            itemExists = true
            existingItem = item
            break
        }
    }

    if itemExists {
        err := db.Q.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
            Quantity: existingItem.Quantity + quantity,
            ID:       existingItem.ID,
            UserID:   user.ID,
        })
        if err != nil {
            log.Printf("Error updating cart item quantity: %v", err)
            return c.String(http.StatusInternalServerError, "Failed to update cart item quantity")
        }
    } else {
        log.Printf("Adding to cart: UserID=%d, MenuItemID=%d, Quantity=%d", user.ID, itemID, 1)
        err := db.Q.AddToCart(ctx, db.AddToCartParams{
            UserID:     user.ID,
            MenuItemID: itemID,
            Quantity:   quantity,
        })
        if err != nil {
            log.Printf("Error adding item to cart: %v", err)
            return c.String(http.StatusInternalServerError, "Failed to add item to cart")
        }
    }

    return c.Redirect(http.StatusSeeOther, "/user/cart")
}

func RemoveAllFromCartHandler(c echo.Context) error {
    user, ok := c.Get("user").(db.User)
    if !ok {
        return c.String(http.StatusInternalServerError, "User not found")
    }
    ctx := c.Request().Context()
    itemID, err := strconv.ParseInt(c.FormValue("item_id"), 10, 64)
    if err != nil {
        return c.String(http.StatusBadRequest, "Invalid item ID")
    }

    cartItems, err := db.Q.GetCartItems(ctx, user.ID)
    if err != nil {
        log.Printf("Error fetching cart items: %v", err)
        return c.String(http.StatusInternalServerError, "Failed to fetch cart items")
    }
    for _, item := range cartItems {
        if itemID == item.MenuItemID {
            err = db.Q.RemoveFromCart(ctx, db.RemoveFromCartParams{
                ID:     item.ID, 
                UserID: user.ID,
            })
            if err != nil {
                log.Printf("Error removing cart item: %v", err)
                return c.String(http.StatusInternalServerError, "Failed to remove cart item")
            }
        }
    }
    return c.Redirect(http.StatusSeeOther, "/user/cart")
}

func RemoveFromCartHandler(c echo.Context) error {
    user, ok := c.Get("user").(db.User)
    if !ok {
        return c.String(http.StatusInternalServerError, "User not found")
    }
    ctx := c.Request().Context()

    itemID, err := strconv.ParseInt(c.FormValue("item_id"), 10, 64)
    if err != nil {
        return c.String(http.StatusBadRequest, "Invalid item ID")
    }

    cartItems, err := db.Q.GetCartItems(ctx, user.ID)
    if err != nil {
        log.Printf("Error fetching cart items: %v", err)
        return c.String(http.StatusInternalServerError, "Failed to fetch cart items")
    }

    itemFound := false
    for _, item := range cartItems {
        if itemID == item.MenuItemID {
            itemFound = true
            if item.Quantity > 1 {
                err := db.Q.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
                    Quantity: item.Quantity - 1,
                    ID:       item.ID,
                    UserID:   user.ID,
                })
                if err != nil {
                    log.Printf("Error updating cart item quantity: %v", err)
                    return c.String(http.StatusInternalServerError, "Failed to update cart item quantity")
                }
            } else {
                err = db.Q.RemoveFromCart(ctx, db.RemoveFromCartParams{
                    ID:     item.ID, 
                    UserID: user.ID,
                })
                if err != nil {
                    log.Printf("Error removing cart item: %v", err)
                    return c.String(http.StatusInternalServerError, "Failed to remove cart item")
                }
            }
            break
        }
    }

    if !itemFound {
        return c.String(http.StatusNotFound, "Item not found in cart")
    }

    return c.Redirect(http.StatusSeeOther, "/user/cart")
}


// func UpdateCartItemHandler(c echo.Context) error {
//     user, ok := c.Get("user").(db.User)
//     if !ok {
//         return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "User not found"})
//     }

//     var request struct {
//         ItemID   int64 `json:"item_id"`
//         Quantity int64 `json:"quantity"`
//     }

//     if err := c.Bind(&request); err != nil {
//         return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "Invalid request"})
//     }

//     ctx := c.Request().Context()
//     err := db.Q.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
//         UserID:     user.ID,
//         ID: request.ItemID,
//         Quantity:   request.Quantity,
//     })

//     if err != nil {
//         log.Printf("Error updating cart item quantity: %v", err)
//         return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "Failed to update cart item"})
//     }

//     return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
// }