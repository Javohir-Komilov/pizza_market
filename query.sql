-- name: CreateUser :one
INSERT INTO users (username, password, email, is_admin)
VALUES (?, ?, ?, ?) RETURNING id;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET email = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: CreateMenuItem :one
INSERT INTO menu_items (name, description, price, image_url, category_id)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: GetMenuItemById :one
SELECT * FROM menu_items WHERE id = ? LIMIT 1;

-- name: GetCategories :many
SELECT * FROM categories;

-- name: GetCategoryById :one
SELECT * FROM categories
WHERE id = ? LIMIT 1;

-- name: ListMenuItems :many
SELECT * FROM menu_items ORDER BY category_id, name;

-- name: UpdateMenuItem :exec
UPDATE menu_items
SET name = ?, description = ?, price = ?, category_id = ?
WHERE id = ?;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = ?;

-- name: CreateOrder :one
INSERT INTO orders (user_id, total_price, status)
VALUES (?, ?, ?) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = ? LIMIT 1;

-- name: ListOrdersByUser :many
SELECT * FROM orders WHERE user_id = ? ORDER BY created_at DESC;

-- name: ListAllOrders :many
SELECT * FROM orders ORDER BY created_at DESC;

-- name: UpdateOrderStatus :exec
UPDATE orders SET status = ? WHERE id = ?;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, menu_item_id, quantity)
VALUES (?, ?, ?) RETURNING *;

-- name: GetOrderItems :many
SELECT oi.*, mi.name, mi.price
FROM order_items oi
JOIN menu_items mi ON oi.menu_item_id = mi.id
WHERE oi.order_id = ?;

-- name: AddToCart :exec 
INSERT INTO cart_items (user_id, menu_item_id, quantity)
VALUES (?, ?, ?);

-- name: GetCartItems :many
SELECT ci.*
FROM cart_items ci
JOIN menu_items mi ON ci.menu_item_id = mi.id
WHERE ci.user_id = ?;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = ? WHERE id = ? AND user_id = ?;

-- name: RemoveFromCart :exec
DELETE FROM cart_items WHERE id = ? AND user_id = ?;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE user_id = ?;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at)
VALUES (?, ?, ?) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE token = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP;

