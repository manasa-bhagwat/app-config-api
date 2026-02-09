package main

import (
	"database/sql"
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

// GetPageByID - GET /pages/:id
func GetPageByID(c *gin.Context) {
	id := c.Param("id") // read ID from URL

	var page Page

	// 1: Fetch page
	err := DB.QueryRow(`
		SELECT id, name, route, is_home, created_at, updated_at
		FROM pages
		WHERE id = $1
	`, id).Scan(
		&page.ID,
		&page.Name,
		&page.Route,
		&page.IsHome,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	// If no page found
	if err == sql.ErrNoRows {
		errorResponse(c, http.StatusNotFound, "NOT_FOUND", "Page not found")
		return
	}

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Database error")
		return
	}

	// 2: Fetch widgets for this page
	rows, err := DB.Query(`
		SELECT id, page_id, type, position, config, created_at, updated_at
		FROM widgets
		WHERE page_id = $1
		ORDER BY position
	`, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch widgets")
		return
	}
	defer rows.Close()

	var widgets []Widget

	for rows.Next() {
		var w Widget
		err := rows.Scan(
			&w.ID,
			&w.PageID,
			&w.Type,
			&w.Position,
			&w.Config,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Failed to read widget data")
			return
		}

		widgets = append(widgets, w)
	}

	// 3: Return page + widgets
	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"widgets": widgets,
	})
}

// UpdatePage - PUT /pages/:id
func UpdatePage(c *gin.Context) {
	id := c.Param("id")

	var input Page

	// 1: Read JSON body
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	// 2: Validate required fields
	if input.Name == "" {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Page name is required")
		return
	}

	if input.Route == "" {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Route is required")
		return
	}

	// 3: Check if page exists
	var existing Page
	err := DB.QueryRow(`
		SELECT id, is_home
		FROM pages
		WHERE id = $1
	`, id).Scan(&existing.ID, &existing.IsHome)

	if err == sql.ErrNoRows {
		errorResponse(c, http.StatusNotFound, "NOT_FOUND", "Page not found")
		return
	}
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Database error")
		return
	}

	// 4: If is_home = true, ensure no other home page exists
	if input.IsHome {
		var count int
		err := DB.QueryRow(`
			SELECT COUNT(*) FROM pages 
			WHERE is_home = true AND id != $1
		`, id).Scan(&count)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Database error")
			return
		}

		if count > 0 {
			errorResponse(c, http.StatusConflict, "CONFLICT", "Another home page already exists")
			return
		}
	}

	// 5: Update page
	query := `
		UPDATE pages
		SET name = $1, route = $2, is_home = $3, updated_at = $4
		WHERE id = $5
	`

	_, err = DB.Exec(query, input.Name, input.Route, input.IsHome, time.Now(), id)
	if err != nil {
		errorResponse(c, http.StatusConflict, "CONFLICT", "Route may already exist")
		return
	}

	// 6: Return updated page
	input.ID = id
	input.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, input)
}

// DeletePage - DELETE /pages/:id
func DeletePage(c *gin.Context) {
	id := c.Param("id")

	// 1: Check if page exists and whether it is home
	var isHome bool
	err := DB.QueryRow(`
		SELECT is_home FROM pages WHERE id = $1
	`, id).Scan(&isHome)

	if err == sql.ErrNoRows {
		errorResponse(c, http.StatusNotFound, "NOT_FOUND", "Page not found")
		return
	}

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Database error")
		return
	}

	// 2: Prevent deleting home page
	if isHome {
		errorResponse(c, http.StatusConflict, "CONFLICT", "Cannot delete the home page")
		return
	}

	// 3: Delete page
	_, err = DB.Exec(`
		DELETE FROM pages WHERE id = $1
	`, id)

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "DB_ERROR", "Failed to delete page")
		return
	}

	// 4: Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Page deleted successfully",
	})
}
