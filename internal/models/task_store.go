package models

type TaskStore interface {
	GetAll() ([]*Task, error)

	Get(id string) (*Task, error)

	Create(task *Task) (*Task, error)

	Update(task *Task) error

	Delete(id string) error
}
