package main

import (
	"encoding/json"
	"time"
)

// Page represents a mobile app screen configuration
type Page struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Route     string    `json:"route"`
	IsHome    bool      `json:"is_home"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Widget represents a UI component placed on a page
type Widget struct {
	ID        string          `json:"id"`
	PageID    string          `json:"page_id"`
	Type      string          `json:"type"`
	Position  int             `json:"position"`
	Config    json.RawMessage `json:"config"` // flexible JSON storage
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ReorderRequest represents new order of widgets
type ReorderRequest struct {
	WidgetIDs []string `json:"widget_ids"`
}
