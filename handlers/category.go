package handlers

import (
	"net/http"
	"pizza/db"
	"pizza/templates"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CategoryHandler(c echo.Context) error {

	ctx := c.Request().Context()
	categoryId, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
	if err != nil {
		return err
	}

	items, err := db.Q.ListMenuItems(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch menu items")
	}
	
	
	category, err := db.Q.GetCategoryById(ctx, categoryId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch categories")
	}
	
	return templates.Render(c, templates.Category(items, category))
}