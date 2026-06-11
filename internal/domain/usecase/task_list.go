package usecase

import "github.com/alexnesterov/task-scheduler/internal/domain/entity"

func (s *TaskService) ListTasks() ([]*entity.Task, error) {
	return s.TaskRepo.List(50)
}
