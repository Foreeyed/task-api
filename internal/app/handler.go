package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler обрабатывает HTTP-запросы для работы с задачами
type Handler struct {
	service Service
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	task, err := h.service.CreateTask(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create task", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID", err)
		return
	}

	task, err := h.service.GetTask(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "task not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get tasks", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID", err)
		return
	}

	if err := h.service.DeleteTask(ctx, id); err != nil {
		respondWithError(w, http.StatusNotFound, "failed to delete task", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CancelTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid task ID", err)
		return
	}

	task, err := h.service.CancelTask(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "task not found":
			status = http.StatusNotFound
		case "only processing tasks can be canceled":
			status = http.StatusConflict
		}
		respondWithError(w, status, "failed to cancel task", err)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// ErrorResponse представляет структуру для ошибок API
type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
}

// respondWithJSON отправляет JSON-ответ
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// respondWithError отправляет JSON-ответ с ошибкой
func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	response := ErrorResponse{
		Timestamp: time.Now(),
		Status:    code,
		Message:   message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	respondWithJSON(w, code, response)
}
