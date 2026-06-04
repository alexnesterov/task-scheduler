// Package httpapi
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
	"github.com/alexnesterov/task-scheduler/internal/domain/port"
)

type TaskHandler struct {
	TaskUseCase port.TaskUseCase
}

func (h *TaskHandler) NextDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowDate, err := parseNow(now)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	nextDate, err := h.TaskUseCase.NextDate(nowDate, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}

func parseNow(n string) (time.Time, error) {
	if n == "" {
		return time.Now(), nil
	}

	now, err := time.Parse("20060102", n)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse now: %v", err)
	}

	return now, nil
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req port.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid request body: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	id, err := h.TaskUseCase.CreateTask(req)
	if errors.Is(err, entity.ErrCreateTask) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := h.TaskUseCase.ListTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
}

func (h *TaskHandler) ReadTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	task, err := h.TaskUseCase.ReadTask(id)
	if errors.Is(err, entity.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req port.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid request body: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	err := h.TaskUseCase.UpdateTask(req)
	if errors.Is(err, entity.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, entity.ErrUpdateTask) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	err := h.TaskUseCase.DeleteTask(id)
	if errors.Is(err, entity.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{})
}

func (h *TaskHandler) DoneTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	err := h.TaskUseCase.DoneTask(id)
	if errors.Is(err, entity.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, entity.ErrDoneTask) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("internal error: %v", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{})
}
