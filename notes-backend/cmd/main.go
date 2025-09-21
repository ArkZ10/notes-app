package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"notes-backend/internal/auth"
	"notes-backend/internal/categories"
	"notes-backend/internal/images"
	"notes-backend/internal/logs"
	"notes-backend/internal/middleware"
	"notes-backend/internal/notes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			body TEXT,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL,
			is_favorite BOOLEAN DEFAULT FALSE,
			visibility TEXT CHECK (visibility IN ('private', 'public', 'shared')) DEFAULT 'private',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS images (
			id SERIAL PRIMARY KEY,
			note_id INT REFERENCES notes(id) ON DELETE CASCADE,
			url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			method TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			request_headers JSONB,
			request_body JSONB,
			response_body JSONB,
			status_code INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return fmt.Errorf("error running migration: %w", err)
		}
	}
	return nil
}

func setupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Apply logging middleware early (so all routes get logged)
	r.Use(middleware.LoggingMiddleware(db))

	// Public routes
	r.Static("/uploads", "./uploads")
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.POST("/register", auth.RegisterHandler(db))
	r.POST("/login", auth.LoginHandler(db))
	r.GET("/logs", logs.ListLogsHandler(db))

	// Protected routes
	notesGroup := r.Group("/notes")
	notesGroup.Use(middleware.JWTMiddleware())
	{
		notesGroup.POST("", notes.CreateNoteHandler(db))
		notesGroup.GET("", notes.ListNotesHandler(db))
		notesGroup.GET("/:id", notes.GetNoteByIDHandler(db))
		notesGroup.PATCH("/:id", notes.UpdateNoteHandler(db))
		notesGroup.DELETE("/:id", notes.DeleteNoteHandler(db))
		notesGroup.POST("/:id/images", images.UploadImageHandler(db))
		notesGroup.GET("/:id/images", images.ListImagesHandler(db))
		notesGroup.DELETE("/:id/images/:image_id", images.DeleteImageHandler(db))
	}

	categoriesGroup := r.Group("/categories")
	categoriesGroup.Use(middleware.JWTMiddleware())
	{
		categoriesGroup.POST("", categories.CreateCategoryHandler(db))
		categoriesGroup.GET("", categories.ListCategoriesHandler(db))
		categoriesGroup.DELETE("/:id", categories.DeleteCategoryHandler(db))
	}

	return r
}

func main() {
	// dsn := "postgres://notesuser:notessecret@localhost:5432/notesdb?sslmode=disable"

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Missing one or more required environment variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	fmt.Println("Connected to PostgreSQL!")

	if err := migrate(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database migrated successfully!")

	r := setupRouter(db)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
