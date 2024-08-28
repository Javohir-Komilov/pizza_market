package main

import (
	"database/sql"
	"log"
	"pizza/db"
	"pizza/handlers"
	"pizza/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DRIVER = "sqlite3"
	DBNAME = "pizza.db"
)

func main() {
	database, err := sql.Open(DRIVER, DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	db.Q = db.New(database)

	e := echo.New()
	e.Debug = true
	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middlewares.Session)

	// Public routes
	e.GET("/", handlers.IndexHandler)
	e.GET("/menu", handlers.MenuHandler)
	e.GET("/category/:id", handlers.CategoryHandler)
	e.Any("/login", handlers.LoginHandler)
	e.Any("/register", handlers.RegisterHandler)
	e.GET("/logout", handlers.LogoutHandler)

	// User routes
	userGroup := e.Group("/user")
	userGroup.Use(middlewares.RequireAuth)
	userGroup.GET("/cart", handlers.CartHandler)
	userGroup.POST("/cart/add", handlers.AddToCartHandler)
	userGroup.POST("/cart/remove", handlers.RemoveFromCartHandler)
	userGroup.POST("/cart/remove-all", handlers.RemoveAllFromCartHandler)
	userGroup.GET("/orders", handlers.UserOrdersHandler)
	userGroup.POST("/place-order", handlers.PlaceOrderHandler)

	// Admin routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middlewares.RequireAuth, middlewares.RequireAdmin)
	adminGroup.Any("/menu/add", handlers.AddMenuItemHandler)
	adminGroup.Any("/menu/edit/:id", handlers.UpdateMenuItemHandler)
	adminGroup.GET("/menu/delete/:id", handlers.DeleteMenuItemHandler)
	adminGroup.GET("/orders", handlers.AdminOrdersHandler)
	adminGroup.POST("/update-status", handlers.UpdateOrderStatusHandler)

	e.Logger.Fatal(e.Start(":8080"))
}