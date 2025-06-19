package cmd

import (
	"log"
	"net/http"

	task "task-api/internal/app"
	"task-api/pkg/server"
)

func main() {
	// Инициализация зависимостей
	taskRepo := task.NewInMemoryRepository()
	taskService := task.NewService(taskRepo)
	taskHandler := task.NewHandler(taskService)

	// Создание HTTP сервера
	server := server.NewServer(taskHandler)

	// Запуск сервера на порту 8080
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
