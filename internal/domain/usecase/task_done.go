package usecase

import (
	"fmt"
	"time"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
)

func (s *TaskService) DoneTask(id string) error {
	task, err := s.TaskRepo.Read(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.TaskRepo.Delete(id)
	}

	now := time.Now()

	next, err := s.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("%w: %w", entity.ErrDoneTask, err)
	}

	task.Date = next

	return s.TaskRepo.Update(task)
}
