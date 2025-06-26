package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/NanobyteRuata/go-taskmanager/internal/models"
	"github.com/google/uuid"
)

type FileStorage struct {
	filename string
	tasks    map[string]*models.Task
	mutex    sync.RWMutex // Read-Write Mutex can be used to prevent multiple goroutines read/write the same file at the same time which can cause inconsistent or corrupted data.
}

func NewFileStorage(filename string) (*FileStorage, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fs := &FileStorage{
		filename: absPath,
		tasks:    make(map[string]*models.Task),
	}

	if err := fs.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	return fs, nil
}

func (fs *FileStorage) save() error {
	tasks := make([]*models.Task, 0, len(fs.tasks))
	for _, task := range fs.tasks {
		tasks = append(tasks, task)
	}

	// os.Create handles both create and truncate
	file, err := os.Create(fs.filename)
	if err != nil {
		return err
	}
	defer file.Close() // defer = close file after function end

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}

func (fs *FileStorage) load() error {
	file, err := os.Open(fs.filename)
	if err != nil {
		return err
	}
	defer file.Close() // defer = close file after function end

	var tasks []*models.Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return err
	}

	fs.tasks = make(map[string]*models.Task, len(tasks))
	for _, task := range tasks {
		fs.tasks[task.ID] = task
	}

	return nil
}

func (fs *FileStorage) GetAll() ([]*models.Task, error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	tasks := make([]*models.Task, 0, len(fs.tasks))
	for _, task := range fs.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (fs *FileStorage) Get(id string) (*models.Task, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	task, exists := fs.tasks[id]
	if !exists {
		return nil, models.ErrTaskNotFound
	}

	return task, nil
}

func (fs *FileStorage) Create(task *models.Task) (*models.Task, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	task.ID = uuid.New().String()

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	fs.tasks[task.ID] = task

	if err := fs.save(); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

func (fs *FileStorage) Update(task *models.Task) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if task.ID == "" {
		return models.ErrInvalidID
	}

	if _, exists := fs.tasks[task.ID]; !exists {
		return models.ErrTaskNotFound
	}

	fs.tasks[task.ID] = task

	if err := fs.save(); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

func (fs *FileStorage) Delete(id string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if _, exists := fs.tasks[id]; !exists {
		return models.ErrTaskNotFound
	}

	delete(fs.tasks, id)
	return fs.save()
}
