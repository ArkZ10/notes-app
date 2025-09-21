package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Bind JSON from request
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Step 2: Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Step 3: Insert user into database
		var id int
		err = db.QueryRow(
			"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
			req.Username, req.Email, string(hash),
		).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Step 4: Respond with success
		c.JSON(http.StatusOK, gin.H{"message": "User registered", "id": id})
	}
}

// LoginHandler handles user login
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Bind JSON
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Step 2: Retrieve user from database
		var id int
		var hash string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", req.Username).Scan(&id, &hash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Step 3: Compare password
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Step 4: Generate JWT using your jwt.go
		tokenString, err := GenerateToken(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Step 5: Respond with token
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
