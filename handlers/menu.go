package handlers

import (
	"net/http"
	"pizza/db"
	"pizza/templates"


	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	return templates.Render(c, templates.Index())
}

func MenuHandler(c echo.Context) error {
	items, err := db.Q.ListMenuItems(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch menu items")
	}
	
	categories, err := db.Q.GetCategories(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch categories")
	}
	return templates.Render(c, templates.Menu(items, categories))
}


