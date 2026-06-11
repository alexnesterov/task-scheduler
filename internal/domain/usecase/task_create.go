package usecase

import (
	"fmt"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
	"github.com/alexnesterov/task-scheduler/internal/domain/port"
)

func (s *TaskService) CreateTask(req port.CreateTaskRequest) (string, error) {
	if req.Title == "" {
		return "", fmt.Errorf("%w: title is required", entity.ErrCreateTask)
	}

	task := &entity.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := s.checkDate(task); err != nil {
		return "", fmt.Errorf("%w: %w", entity.ErrCreateTask, err)
	}

	id, err := s.TaskRepo.Create(task)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(id), nil
}
