package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexnesterov/task-scheduler/internal/adapter/httpapi"
	"github.com/alexnesterov/task-scheduler/internal/db"
	"github.com/alexnesterov/task-scheduler/internal/domain/usecase"
	"github.com/alexnesterov/task-scheduler/internal/infrastructure/sqlite"
)

const PORT = "7540"

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer func() { _ = db.DB.Close() }()

	router := http.NewServeMux()

	taskRepository := sqlite.NewTaskRepository(db.DB)
	taskService := &usecase.TaskUseCase{
		TaskRepo: taskRepository,
	}

	handler := httpapi.TaskHandler{
		TaskUseCase: taskService,
	}

	router.HandleFunc("/", http.FileServer(http.Dir("web")).ServeHTTP)
	router.HandleFunc("GET /api/nextdate", handler.NextDate)
	router.HandleFunc("POST /api/task", handler.CreateTask)
	router.HandleFunc("GET /api/tasks", handler.ListTasks)
	router.HandleFunc("GET /api/task", handler.ReadTask)
	router.HandleFunc("PUT /api/task", handler.UpdateTask)
	router.HandleFunc("DELETE /api/task", handler.DeleteTask)
	router.HandleFunc("POST /api/task/done", handler.DoneTask)

	server := http.Server{
		Addr:         ":" + PORT,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Server is running on port %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
