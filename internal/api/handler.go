package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/just4fun-xd/task-manager/internal/task"
)

type Handler struct {
	service *task.Service
}

func NewHandler(s *task.Service) *Handler {
	return &Handler{
		service: s,
	}
}

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupID     *int   `json:"group_id"`
}

type UpdateTaskRequest struct {
	CreateTaskRequest
	Status task.TaskStatus `json:"status"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	t, err := h.service.CreateTask(r.Context(), req.Name, req.Description, req.GroupID)
	if err != nil {
		if errors.Is(err, task.ErrEmptyTaskName) || errors.Is(err, task.ErrGroupNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok := h.getId(w, r)
	if !ok {
		return
	}
	t, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	t, err := h.service.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := h.getId(w, r)
	if !ok {
		return
	}
	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, ok := h.getId(w, r)
	if !ok {
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if !req.Status.IsValid() {
		http.Error(w, "invalid task status", http.StatusBadRequest)
		return
	}
	t, err := h.service.UpdateTask(r.Context(), id, req.Name, req.Description, req.Status)
	if err != nil {
		if errors.Is(err, task.ErrEmptyTaskName) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, task.ErrNewTaskStatus) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, task.ErrDoneEdit) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) getId(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
