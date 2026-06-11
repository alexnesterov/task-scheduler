package usecase

import (
	"fmt"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
	"github.com/alexnesterov/task-scheduler/internal/domain/port"
)

func (s *TaskService) UpdateTask(req port.UpdateTaskRequest) error {
	if req.Title == "" {
		return fmt.Errorf("%w: title is required", entity.ErrUpdateTask)
	}

	task := &entity.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := s.checkDate(task); err != nil {
		return fmt.Errorf("%w: %w", entity.ErrUpdateTask, err)
	}

	return s.TaskRepo.Update(task)
}
