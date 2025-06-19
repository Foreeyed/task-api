package app

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Repository определяет интерфейс для работы с задачами
type Repository interface {
	Create(task Task) (*Task, error)
	GetByID(id string) (*Task, error)
	GetAll() ([]Task, error)
	Delete(id string) error
	Update(task Task) (*Task, error)
}

// inMemoryRepository - реализация Repository, хранящая данные в памяти
type inMemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

// NewInMemoryRepository создает новый экземпляр inMemoryRepository
func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		tasks: make(map[string]Task),
	}
}

// Create добавляет новую задачу в хранилище
func (r *inMemoryRepository) Create(task Task) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем, не существует ли уже задача с таким ID
	if _, exists := r.tasks[task.ID]; exists {
		return nil, errors.New("task with this ID already exists")
	}

	// Устанавливаем время создания, если оно не задано
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	r.tasks[task.ID] = task
	return &task, nil
}

// GetByID возвращает задачу по её ID
func (r *inMemoryRepository) GetByID(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return &task, nil
}

// GetAll возвращает все задачи из хранилища
func (r *inMemoryRepository) GetAll() ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Delete удаляет задачу по её ID
func (r *inMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)
	return nil
}

// Update обновляет существующую задачу
func (r *inMemoryRepository) Update(task Task) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем, существует ли задача
	if _, exists := r.tasks[task.ID]; !exists {
		return nil, errors.New("task not found")
	}

	// Обновляем задачу
	r.tasks[task.ID] = task
	return &task, nil
}

func generateID() string {
	return uuid.NewString()
}
