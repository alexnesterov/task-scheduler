package main

import (
	"log"
	"net/http"

	"github.com/alexnesterov/task-scheduler/internal/adapter/httpapi"
	"github.com/alexnesterov/task-scheduler/internal/config"
	"github.com/alexnesterov/task-scheduler/internal/db"
	"github.com/alexnesterov/task-scheduler/internal/domain/usecase"
	"github.com/alexnesterov/task-scheduler/internal/infrastructure/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = db.Init(cfg.DB.File)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer func() { _ = db.DB.Close() }()

	router := http.NewServeMux()

	taskRepository := sqlite.NewTaskRepository(db.DB)
	taskService := &usecase.TaskService{
		TaskRepo: taskRepository,
	}

	handler := httpapi.TaskHandler{
		TaskUseCase: taskService,
	}

	router.HandleFunc("GET /api/nextdate", handler.NextDate)
	router.HandleFunc("POST /api/task", handler.CreateTask)
	router.HandleFunc("GET /api/tasks", handler.ListTasks)
	router.HandleFunc("GET /api/task", handler.ReadTask)
	router.HandleFunc("PUT /api/task", handler.UpdateTask)
	router.HandleFunc("DELETE /api/task", handler.DeleteTask)
	router.HandleFunc("POST /api/task/done", handler.DoneTask)

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	log.Printf("Server is running on port %s", cfg.HTTP.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
