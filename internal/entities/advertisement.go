// internal/entities/advertisement.go
package entities

import "time"

type Advertisement struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    CreatedAt   time.Time `json:"created_at"`
    IsActive    bool      `json:"is_active"`
}
