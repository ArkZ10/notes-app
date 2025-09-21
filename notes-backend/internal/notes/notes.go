package notes

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct for creating a new note
type NoteRequest struct {
	Title      string `json:"title" binding:"required"`
	Body       string `json:"body"`
	CategoryID *int   `json:"category_id"` // optional
	IsFavorite bool   `json:"is_favorite"`
	Visibility string `json:"visibility"` // "private", "public", "shared"
}

// Struct for listing/updating notes
type Note struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	CategoryID *int      `json:"category_id,omitempty"`
	IsFavorite bool      `json:"is_favorite"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func CreateNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req NoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userID := c.GetInt("userID")

		if req.Visibility == "" {
			req.Visibility = "private"
		}

		now := time.Now()

		var noteID int
		err := db.QueryRow(`
			INSERT INTO notes (user_id, title, body, category_id, is_favorite, visibility, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, userID, req.Title, req.Body, req.CategoryID, req.IsFavorite, req.Visibility, now, now).Scan(&noteID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Note created successfully",
			"note_id": noteID,
		})
	}
}

func ListNotesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")

		// Read query params
		favorite := c.Query("favorite")
		search := c.Query("search")
		categoryID := c.Query("category_id")

		// Base query
		query := `
			SELECT id, title, body, category_id, is_favorite, visibility, created_at, updated_at
			FROM notes
			WHERE user_id = $1 AND deleted_at IS NULL
		`
		args := []interface{}{userID}
		i := 2

		// Filter: favorites
		if favorite == "true" {
			query += fmt.Sprintf(" AND is_favorite = $%d", i)
			args = append(args, true)
			i++
		}

		// Filter: search
		if search != "" {
			query += fmt.Sprintf(" AND (title ILIKE $%d OR body ILIKE $%d)", i, i+1)
			args = append(args, "%"+search+"%", "%"+search+"%")
			i += 2
		}

		// Filter: category
		if categoryID != "" {
			if categoryID == "none" {
				query += " AND category_id IS NULL"
			} else {
				query += fmt.Sprintf(" AND category_id = $%d", i)
				args = append(args, categoryID)
				i++
			}
		}

		// Order
		query += " ORDER BY created_at DESC"

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
			return
		}
		defer rows.Close()

		var notes []Note
		for rows.Next() {
			var n Note
			var category sql.NullInt64

			if err := rows.Scan(&n.ID, &n.Title, &n.Body, &category, &n.IsFavorite, &n.Visibility, &n.CreatedAt, &n.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan note"})
				return
			}

			if category.Valid {
				id := int(category.Int64)
				n.CategoryID = &id
			}

			notes = append(notes, n)
		}

		c.JSON(http.StatusOK, gin.H{"notes": notes})
	}
}

type UpdateNoteRequest struct {
	Title      *string `json:"title"`
	Body       *string `json:"body"`
	CategoryID *int    `json:"category_id"`
	IsFavorite *bool   `json:"is_favorite"`
	Visibility *string `json:"visibility"`
}

func UpdateNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")

		var req UpdateNoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Build dynamic update query
		query := "UPDATE notes SET "
		args := []interface{}{}
		i := 1

		if req.Title != nil {
			query += fmt.Sprintf("title=$%d,", i)
			args = append(args, *req.Title)
			i++
		}
		if req.Body != nil {
			query += fmt.Sprintf("body=$%d,", i)
			args = append(args, *req.Body)
			i++
		}
		if req.CategoryID != nil {
			query += fmt.Sprintf("category_id=$%d,", i)
			args = append(args, *req.CategoryID)
			i++
		}
		if req.IsFavorite != nil {
			query += fmt.Sprintf("is_favorite=$%d,", i)
			args = append(args, *req.IsFavorite)
			i++
		}
		if req.Visibility != nil {
			query += fmt.Sprintf("visibility=$%d,", i)
			args = append(args, *req.Visibility)
			i++
		}

		if len(args) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		// Update updated_at timestamp
		query += fmt.Sprintf("updated_at=$%d", i)
		args = append(args, time.Now())
		i++

		query += fmt.Sprintf(" WHERE id=$%d AND user_id=$%d AND deleted_at IS NULL", i, i+1)
		args = append(args, noteID, userID)

		res, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
	}
}

func DeleteNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")

		res, err := db.Exec(`
			UPDATE notes
			SET deleted_at=$1
			WHERE id=$2 AND user_id=$3 AND deleted_at IS NULL
		`, time.Now(), noteID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
	}
}

func GetNoteByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		noteID := c.Param("id")

		var n Note
		var categoryID sql.NullInt64

		query := `
            SELECT id, title, body, category_id, is_favorite, visibility, created_at, updated_at
            FROM notes
            WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
        `
		err := db.QueryRow(query, noteID, userID).Scan(
			&n.ID, &n.Title, &n.Body, &categoryID,
			&n.IsFavorite, &n.Visibility, &n.CreatedAt, &n.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note"})
			return
		}

		if categoryID.Valid {
			id := int(categoryID.Int64)
			n.CategoryID = &id
		}

		c.JSON(http.StatusOK, n)
	}
}
