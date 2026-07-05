package models

import "time"

type Task struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	Priority string `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Priority string `json:"priority"`
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
	Description *string `json:"description"`
	Priority *string `json:"priority"`
	Status *string `json:"status"`
}
