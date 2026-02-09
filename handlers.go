package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// errorResponse - for a consistent error format
// example - errorResponse(c, 500, "DB_ERROR", "Database error")
func errorResponse(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// CreatePage - POST /pages
func CreatePage(c *gin.Context) {
	var page Page

	// 1. Read JSON request body into Page struct
	if err := c.ShouldBindJSON(&page); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	// 2. Basic Validations
	if page.Name == "" {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Page name is required")
		return
	}

	if page.Route == "" {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Route is required")
		return
	}

	// 3. If is_home = true, check if another home page exists
	if page.IsHome {
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM pages WHERE is_home = true").Scan(&count)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Database error")
			return
		}

		if count > 0 {
			errorResponse(c, http.StatusConflict, "CONFLICT", "Home page already exists")
			return
		}
	}

	// 4. Generate UUID for new page
	page.ID = uuid.New().String()
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	// 5. Insert into db
	query := `INSERT INTO pages (id, name, route, is_home, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := DB.Exec(query, page.ID, page.Name, page.Route, page.IsHome, page.CreatedAt, page.UpdatedAt)
	if err != nil {
		errorResponse(c, http.StatusConflict, "CONFLICT", "Page route already exists")
		return
	}

	// 6. Return created page
	c.JSON(http.StatusCreated, page)
}

// GetPages - GET /pages
func GetPages(c *gin.Context) {
	var pages []Page

	// Query all pages
	rows, err := DB.Query(`
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		ORDER BY created_at
	`)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch pages")
		return
	}

	defer rows.Close()

	// Iterate over result set
	for rows.Next() {
		var page Page
		err := rows.Scan(
			&page.ID,
			&page.Name,
			&page.Route,
			&page.IsHome,
			&page.CreatedAt,
			&page.UpdatedAt,
		)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Failed to read page data")
			return
		}

		pages = append(pages, page)
	}

	// Return list of pages
	c.JSON(http.StatusOK, pages)

}
