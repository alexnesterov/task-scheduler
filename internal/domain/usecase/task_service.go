// Package usecase
package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
	"github.com/alexnesterov/task-scheduler/internal/domain/port"
)

type TaskService struct {
	TaskRepo port.TaskRepository
}

func (s *TaskService) NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("parse date: %w", err)
	}

	repeated := strings.Split(repeat, " ")

	switch repeated[0] {
	case "d":
		if len(repeated) < 2 {
			return "", fmt.Errorf("empty interval")
		}

		interval, err := strconv.Atoi(repeated[1])
		if err != nil {
			return "", fmt.Errorf("parse interval: %w", err)
		}

		if interval <= 0 {
			return "", fmt.Errorf("interval must be positive")
		}

		if interval > 400 {
			return "", fmt.Errorf("max days 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if s.afterNow(date, now) {
				break
			}
		}

	case "y":
		interval := 1

		for {
			date = date.AddDate(interval, 0, 0)
			if s.afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unknown repeat format: %v", repeat)
	}

	return date.Format("20060102"), nil
}

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

func (s *TaskService) ListTasks() ([]*entity.Task, error) {
	return s.TaskRepo.List(50)
}

func (s *TaskService) ReadTask(id string) (*entity.Task, error) {
	return s.TaskRepo.Read(id)
}

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

func (s *TaskService) DeleteTask(id string) error {
	return s.TaskRepo.Delete(id)
}

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

func (s *TaskService) checkDate(task *entity.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(`20060102`)
	}

	t, err := time.Parse(`20060102`, task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = s.NextDate(now, task.Date, task.Repeat)
	}
	if err != nil {
		return err
	}

	if s.afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(`20060102`)
		} else {
			task.Date = next
		}
	}

	return nil
}

func (s *TaskService) afterNow(date, now time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}
