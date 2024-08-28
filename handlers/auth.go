package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"pizza/db"
	"pizza/templates"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	if c.Request().Method == "POST" {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if username == "" || password == "" {
			return c.String(400, "username or password is empty")
		}

		ctx := c.Request().Context()

		// Get the expected password from our in memory map
		user, err := db.Q.GetUserByUsername(ctx, username)

		hashStr := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

		if err != nil || user.Password != hashStr {
			return c.String(http.StatusUnauthorized, "wrong username or password")
		}

		// Create a new random session token
		// we use the "github.com/google/uuid" library to generate UUIDs
		err = createSession(user.ID, c)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to create session")
		}

		return c.Redirect(http.StatusFound, "/menu ")

	}

	return templates.Render(c, templates.Login())
}

func RegisterHandler(c echo.Context) error {
	if c.Request().Method == "POST" {
		username := strings.TrimSpace(c.FormValue("username"))
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirm_password")
		email := strings.TrimSpace(c.FormValue("email"))

		if username == "" || password == "" || email == "" {
			return c.String(http.StatusBadRequest, "All fields are required")
		}

		if password != confirmPassword {
			return c.String(http.StatusBadRequest, "Passwords do not match")
		}
		
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(email) {
			return c.String(http.StatusBadRequest, "Invalid email format")
		}

		ctx := c.Request().Context()

		if user, err := db.Q.GetUserByUsername(ctx, username); err == nil && user.Username == username {
			return c.String(http.StatusBadRequest, "Username already exists")
		}

		if user, err := db.Q.GetUserByEmail(ctx, email); err == nil && user.Email == email {
			return c.String(http.StatusBadRequest, "Email already in use")
		}

		hash := sha256.Sum256([]byte(password))
		hashStr := fmt.Sprintf("%x", hash)


		userId, err := db.Q.CreateUser(ctx, db.CreateUserParams{
			Username: username,
			Password: hashStr,
			Email:    email,
			IsAdmin:  int64(0),
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create user")
		}

		err = createSession(userId, c)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create session")
		}

		return c.Redirect(http.StatusFound, "/menu")
	}

	return templates.Render(c, templates.Register())
}

func LogoutHandler(c echo.Context) error {
	// Delete the session
	sessionCookie, _ := c.Cookie("session_token")
	db.Q.DeleteSession(c.Request().Context(), sessionCookie.Value)
	// Redirect to the main page
	return c.Redirect(http.StatusFound, "/")
}

func createSession(userId int64, c echo.Context) error {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Minute)

	_, err := db.Q.CreateSession(c.Request().Context(), db.CreateSessionParams{
		UserID:    userId,
		Token:     sessionToken,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		return err
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	c.SetCookie(&http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	return nil
}
