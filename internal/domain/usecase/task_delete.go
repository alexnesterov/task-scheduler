package usecase

func (s *TaskService) DeleteTask(id string) error {
	return s.TaskRepo.Delete(id)
}
