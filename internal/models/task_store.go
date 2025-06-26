package models

type TaskStore interface {
	GetAll() ([]*Task, error)

	Create(task *Task) (*Task, error)
}
