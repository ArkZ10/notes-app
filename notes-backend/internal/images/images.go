package images

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Image struct {
	ID        int    `json:"id"`
	NoteID    int    `json:"note_id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

// ✅ Upload Image for a Note
func UploadImageHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")

		// Check if note belongs to user
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL)
		`, noteID, userID).Scan(&exists)

		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not owned by user"})
			return
		}

		// Parse file from form-data (key = "image")
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			return
		}

		// Save file locally (uploads folder)
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		baseURL := "http://localhost:8080/" // later replace with domain
		fileURL := baseURL + filename

		// Insert into images table
		var imageID int
		err = db.QueryRow(`
			INSERT INTO images (note_id, url, created_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`, noteID, fileURL, time.Now()).Scan(&imageID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image metadata"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Image uploaded successfully",
			"image_id": imageID,
			"url":      fileURL,
		})
	}
}

func ListImagesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")

		// Check if note belongs to user
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL)
		`, noteID, userID).Scan(&exists)

		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not owned by user"})
			return
		}

		rows, err := db.Query(`
			SELECT id, url, created_at
			FROM images
			WHERE note_id=$1
			ORDER BY created_at DESC
		`, noteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
			return
		}
		defer rows.Close()

		type Image struct {
			ID        int       `json:"id"`
			URL       string    `json:"url"`
			CreatedAt time.Time `json:"created_at"`
		}

		var images []Image
		for rows.Next() {
			var img Image
			if err := rows.Scan(&img.ID, &img.URL, &img.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning image"})
				return
			}
			images = append(images, img)
		}

		c.JSON(http.StatusOK, gin.H{"images": images})
	}
}

func DeleteImageHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")
		imageID := c.Param("image_id")

		// Check if note belongs to the user
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL)
		`, noteID, userID).Scan(&exists)
		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not owned by user"})
			return
		}

		// Get image path from DB
		var imagePath string
		err = db.QueryRow(`
			SELECT url FROM images WHERE id=$1 AND note_id=$2
		`, imageID, noteID).Scan(&imagePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		// Delete DB record
		_, err = db.Exec(`DELETE FROM images WHERE id=$1 AND note_id=$2`, imageID, noteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
			return
		}

		// Delete file from uploads folder
		if err := os.Remove(imagePath); err != nil {
			// not fatal, but log it
			c.JSON(http.StatusOK, gin.H{
				"message": "Image record deleted, but file could not be removed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
	}
}
