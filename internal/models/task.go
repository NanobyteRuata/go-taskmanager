package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	DueDate     time.Time `json:"due_date"`
}

func NewTask(title string) *Task {
	return &Task{
		Title:     title,
		CreatedAt: time.Now(),
		Completed: false,
	}
}
