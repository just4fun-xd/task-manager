package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/just4fun-xd/task-manager/internal/task"
)

type GroupHandler struct {
	service *task.Service
}

func NewGroupHandler(service *task.Service) *GroupHandler {
	return &GroupHandler{
		service: service,
	}
}

type GroupRequest struct {
	Name string `json:"name"`
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid requst body", http.StatusBadRequest)
		return
	}
	g, err := h.service.CreateGroup(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, task.ErrEmptyGroupName) || errors.Is(err, task.ErrNotUniqGroup) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := GetId(w, r)
	if !ok {
		return
	}

	var req GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	g, err := h.service.UpdateGroup(r.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, task.ErrEmptyGroupName) || errors.Is(err, task.ErrNotUniqGroup) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, task.ErrGroupNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(g)
}
