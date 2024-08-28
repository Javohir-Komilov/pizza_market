package models

import "pizza/db"

type CartItemInfo struct {
    MenuItem db.MenuItem
    Quantity int64
    Subtotal float64
}