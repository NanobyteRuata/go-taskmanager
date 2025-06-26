package models

import (
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
)

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

func (t *Task) Complete() {
	t.Completed = true
	t.CompletedAt = time.Now()
}

func (t *Task) IsOverdue() bool {
	return !t.Completed && !t.DueDate.IsZero() && time.Now().After(t.DueDate)
}
