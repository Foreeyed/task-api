package server

import (
	"net/http"
	task "task-api/internal/app"


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router      *chi.Mux
	taskHandler *task.Handler
}

func NewServer(taskHandler *task.Handler) *Server {
	s := &Server{
		router:      chi.NewRouter(),
		taskHandler: taskHandler,
	}

	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", s.taskHandler.CreateTask)
		r.Get("/", s.taskHandler.GetAllTasks)
		r.Get("/{id}", s.taskHandler.GetTask)
		r.Delete("/{id}", s.taskHandler.DeleteTask)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
