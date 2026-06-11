package usecase

import "github.com/alexnesterov/task-scheduler/internal/domain/entity"

func (s *TaskService) ReadTask(id string) (*entity.Task, error) {
	return s.TaskRepo.Read(id)
}
