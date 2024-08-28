package middlewares

import (
	"pizza/db"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()
		if path == "/login" || path == "/register" {
			return next(c)
		}

		tokenCookie, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, continue without a session
				return next(c)
			}
			// For any other type of error, log it and continue
			c.Logger().Error("Error retrieving session cookie:", err)
			return next(c)
		}

		token := tokenCookie.Value
		ctx := c.Request().Context()
		userSession, err := db.Q.GetSessionByToken(ctx, token)

		if err != nil {
			// If the session token is not found, continue without a session
			return next(c)
		}

		// Check if the session has expired
		if userSession.ExpiresAt.Before(time.Now()) {
			// Delete the expired session
			if err := db.Q.DeleteSession(ctx, token); err != nil {
				c.Logger().Error("Error deleting expired session:", err)
			}
			// Continue without a session
			return next(c)
		}

		user, err := db.Q.GetUserByID(ctx, userSession.UserID)
		if err != nil {
			c.Logger().Error("Failed to get user:", err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		// Set user in the context
		ctx = context.WithValue(ctx, "user", user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

// RequireAuth middleware ensures that the user is authenticated
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user, ok := ctx.Value("user").(db.User)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		// Set the user in the echo.Context for easy access in handlers
		c.Set("user", user)
		return next(c)
	}
}

// RequireAdmin middleware ensures that the user is authenticated and has admin privileges
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user, ok := ctx.Value("user").(db.User)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		
		// Check if the user is an admin
		if user.IsAdmin == 0 {
			return c.String(http.StatusForbidden, "Admin access required")
		}
		
		// Set the user in the echo.Context for easy access in handlers
		c.Set("user", user)
		return next(c)
	}
}