package tasks

import "time"

// Task is the domain model for a to-do item.
// This is the shape used throughout your service and exposed via HTTP/GraphQL.
type Task struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
