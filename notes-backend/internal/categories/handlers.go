package categories

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// CreateCategoryHandler - POST /categories
func CreateCategoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userID := c.GetInt("userID") // from JWT middleware

		var id int
		err := db.QueryRow(
			"INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id",
			userID, input.Name,
		).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id, "name": input.Name})
	}
}

// ListCategoriesHandler - GET /categories
func ListCategoriesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")

		rows, err := db.Query("SELECT id, name, created_at FROM categories WHERE user_id=$1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}
		defer rows.Close()

		var categories []Category
		for rows.Next() {
			var cat Category
			cat.UserID = userID
			if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt); err != nil {
				continue
			}
			categories = append(categories, cat)
		}

		c.JSON(http.StatusOK, categories)
	}
}

// DeleteCategoryHandler - DELETE /categories/:id
func DeleteCategoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		categoryID := c.Param("id")

		// Option 1: Soft-delete notes? Or set category_id to null
		_, err := db.Exec(`
            UPDATE notes
            SET category_id = NULL
            WHERE category_id = $1 AND user_id = $2
        `, categoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notes"})
			return
		}

		// Delete the category itself
		res, err := db.Exec(`
            DELETE FROM categories
            WHERE id=$1 AND user_id=$2
        `, categoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or not owned by user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}
