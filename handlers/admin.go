package handlers

import (
	"database/sql"
	"net/http"
	"pizza/db"
	"pizza/templates"
	"strconv"

	"github.com/labstack/echo/v4"
)
func AddMenuItemHandler(c echo.Context) error {
	ctx := c.Request().Context()

	if c.Request().Method == "POST" {
		title := c.FormValue("title")
		description := c.FormValue("description")
		price, err := strconv.ParseFloat(c.FormValue("price"), 64)

		if err != nil {
            return err
        }

		image := c.FormValue("image")
		
		categoryId, err := strconv.Atoi(c.FormValue("category"))

		if err != nil {
			return err
		}
		
		catId := sql.NullInt64{Int64: int64(categoryId), Valid: true}
		if categoryId == 0 {
			catId = sql.NullInt64{Valid: false}
		}

		if title == "" || image == "" {
			return c.String(http.StatusBadRequest, "title, image or price is empty")
		}

		db.Q.CreateMenuItem(ctx, db.CreateMenuItemParams{
			Name:          title,
			Description:   sql.NullString{String: description, Valid: true},
			Price:         price,
			ImageUrl:      image,
			CategoryID:    catId,
		})

		return c.Redirect(http.StatusFound, "/menu")
	}

	categories, err := db.Q.GetCategories(ctx)
	if err != nil {
		return err
	}

	return templates.Render(c, templates.Add(categories))
}

func UpdateMenuItemHandler(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	menuItem, err := db.Q.GetMenuItemById(ctx, int64(id))
	if err != nil {
		return err
	}
	
	var category db.Category
	if menuItem.CategoryID.Valid {
		category, err = db.Q.GetCategoryById(ctx, menuItem.CategoryID.Int64)
		if err != nil {
			return err
		}
	}
	categories, err := db.Q.GetCategories(ctx)
	if err != nil {
		return err
	}
	
	if c.Request().Method == "POST" {
		title := c.FormValue("title")
		description := c.FormValue("description")
		price, err := strconv.ParseFloat(c.FormValue("price"), 64)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch price")
		} 
		image := c.FormValue("image")
		categoryId, err := strconv.Atoi(c.FormValue("category"))
		if err != nil {
			return err
		}
		catId := sql.NullInt64{Int64: int64(categoryId), Valid: true}

		if categoryId == 0 {
			catId = sql.NullInt64{Valid: false}
		}

		if title == "" || image == "" {
			return err
		}

		db.Q.UpdateMenuItem(ctx, db.UpdateMenuItemParams{
			Name:        title,
			Description: sql.NullString{String: description, Valid: true},
			Price:       price,
			CategoryID:  catId,
			ID:          int64(id),
		})

		return c.Redirect(http.StatusFound, "/menu")
	}
	return templates.Render(c, templates.Update(categories, category, menuItem))
}

func DeleteMenuItemHandler(c echo.Context) error {
	ctx := c.Request().Context()

    id, err := strconv.Atoi(c.Param("id"))
    if err!= nil {
        return err
    }

	err = db.Q.DeleteMenuItem(ctx, int64(id))

	if err!= nil {
        return c.String(http.StatusInternalServerError, "Failed to delete menu item")
    }

    return c.Redirect(http.StatusSeeOther, "/menu") 
}

func AdminOrdersHandler(c echo.Context) error {
	orders, err := db.Q.ListAllOrders(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch orders")
	}
	// Assume you have an admin orders template
	return templates.Render(c, templates.AdminOrders(orders))
}

// UpdateOrderStatusHandler handles updating the status of an order
func UpdateOrderStatusHandler(c echo.Context) error {
	orderID, _ := strconv.Atoi(c.FormValue("order_id"))
	status := c.FormValue("status")

	err := db.Q.UpdateOrderStatus(c.Request().Context(), db.UpdateOrderStatusParams{
		ID:     int64(orderID),
		Status: status,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update order status")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/orders")
}