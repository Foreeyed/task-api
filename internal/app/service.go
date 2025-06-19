package app

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Service реализует бизнес-логику работы с задачами
type Service interface {
	CreateTask(ctx context.Context) (*Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	DeleteTask(ctx context.Context, id string) error
	CancelTask(ctx context.Context, id string) (*Task, error) // Дополнительный метод для отмены задачи
}

type service struct {
	repo Repository
}

// NewService создает новый экземпляр сервиса
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateTask создает новую задачу и запускает ее обработку в фоне
func (s *service) CreateTask(ctx context.Context) (*Task, error) {
	task := Task{
		ID:        uuid.NewString(), // Генерация UUID
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	// Сохраняем задачу в репозитории
	createdTask, err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	// Запускаем обработку задачи в отдельной горутине
	go s.processTask(ctx, createdTask.ID)

	return createdTask, nil
}

// processTask имитирует долгую I/O операцию (3-5 минут)
func (s *service) processTask(ctx context.Context, id string) {
	// Обновляем статус задачи на "processing"
	_, err := s.getAndUpdateTask(id, func(t *Task) {
		t.Status = StatusProcessing
		now := time.Now()
		t.StartedAt = &now
	})
	if err != nil {
		return
	}

	// Имитация долгой операции с проверкой контекста
	select {
	case <-time.After(time.Minute * time.Duration(3+rand.Intn(3))): // Случайное время 3-5 минут
		// Успешное завершение
		s.completeTask(id, "Task completed successfully", nil)
	case <-ctx.Done():
		// Задача была отменена
		s.completeTask(id, "", errors.New("task canceled"))
	}
}

// completeTask завершает задачу с указанным результатом или ошибкой
func (s *service) completeTask(id string, result string, err error) {
	_, _ = s.getAndUpdateTask(id, func(t *Task) {
		now := time.Now()
		t.CompletedAt = &now
		if t.StartedAt != nil {
			t.Duration = now.Sub(*t.StartedAt)
		}

		if err != nil {
			t.Status = StatusFailed
			errMsg := err.Error()
			t.Error = &errMsg
		} else {
			t.Status = StatusCompleted
			t.Result = &result
		}
	})
}

// GetTask возвращает задачу по ID
func (s *service) GetTask(ctx context.Context, id string) (*Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// GetAllTasks возвращает все задачи
func (s *service) GetAllTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("failed to get tasks")
	}
	return tasks, nil
}

// DeleteTask удаляет задачу по ID
func (s *service) DeleteTask(ctx context.Context, id string) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("failed to delete task")
	}
	return nil
}

// CancelTask отменяет выполнение задачи
func (s *service) CancelTask(ctx context.Context, id string) (*Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.Status != StatusProcessing {
		return nil, errors.New("only processing tasks can be canceled")
	}

	// В реальном приложении здесь должна быть логика отмены
	// Для примера просто помечаем как отмененную
	updatedTask, err := s.getAndUpdateTask(id, func(t *Task) {
		t.Status = StatusFailed
		errMsg := "task canceled by user"
		t.Error = &errMsg
	})
	if err != nil {
		return nil, errors.New("failed to cancel task")
	}

	return updatedTask, nil
}

// getAndUpdateTask вспомогательный метод для атомарного обновления задачи
func (s *service) getAndUpdateTask(id string, updateFunc func(*Task)) (*Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateFunc(task)

	updatedTask, err := s.repo.Update(*task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}
