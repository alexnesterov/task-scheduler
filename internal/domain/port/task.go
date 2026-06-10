// Package port
package port

import (
	"time"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
)

type TaskRepository interface {
	Create(task *entity.Task) (int64, error)
	Read(id string) (*entity.Task, error)
	List(limit int) ([]*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id string) error
}

type CreateTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type UpdateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskService interface {
	NextDate(now time.Time, dstart, repeat string) (string, error)
	CreateTask(req CreateTaskRequest) (string, error)
	ListTasks() ([]*entity.Task, error)
	ReadTask(id string) (*entity.Task, error)
	UpdateTask(req UpdateTaskRequest) error
	DeleteTask(id string) error
	DoneTask(id string) error
}
