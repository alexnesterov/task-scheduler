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

type TaskUseCase struct {
	TaskRepo port.TaskRepository
}

func (uc *TaskUseCase) CreateTask(req port.CreateTaskRequest) (string, error) {
	if req.Title == "" {
		return "", fmt.Errorf("title is required")
	}

	task := &entity.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := uc.checkDate(task); err != nil {
		return "", err
	}

	id, err := uc.TaskRepo.CreateTask(task)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(id), nil
}

func (uc *TaskUseCase) ListTasks() ([]*entity.Task, error) {
	return uc.TaskRepo.ListTasks(50)
}

func (uc *TaskUseCase) ReadTask(id string) (*entity.Task, error) {
	return uc.TaskRepo.ReadTask(id)
}

func (uc *TaskUseCase) UpdateTask(req port.UpdateTaskRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	task := &entity.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := uc.checkDate(task); err != nil {
		return err
	}

	return uc.TaskRepo.UpdateTask(task)
}

func (uc *TaskUseCase) DeleteTask(id string) error {
	return uc.TaskRepo.DeleteTask(id)
}

func (uc *TaskUseCase) DoneTask(id string) error {
	task, err := uc.TaskRepo.ReadTask(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return uc.TaskRepo.DeleteTask(id)
	}

	now := time.Now()

	next, err := uc.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = next

	return uc.TaskRepo.UpdateTask(task)
}

func (uc *TaskUseCase) NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("parse date: %v", err)
	}

	repeated := strings.Split(repeat, " ")

	switch repeated[0] {
	case "d":
		if len(repeated) < 2 {
			return "", fmt.Errorf("empty interval")
		}

		interval, err := strconv.Atoi(repeated[1])
		if err != nil {
			return "", fmt.Errorf("parse interval: %v", err)
		}

		if interval <= 0 {
			return "", fmt.Errorf("interval must be positive")
		}

		if interval > 400 {
			return "", fmt.Errorf("max days 400")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if uc.afterNow(date, now) {
				break
			}
		}

	case "y":
		interval := 1

		for {
			date = date.AddDate(interval, 0, 0)
			if uc.afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unknown repeat format: %v", repeat)
	}

	return date.Format("20060102"), nil
}

func (uc *TaskUseCase) checkDate(task *entity.Task) error {
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
		next, err = uc.NextDate(now, task.Date, task.Repeat)
	}
	if err != nil {
		return err
	}

	if uc.afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(`20060102`)
		} else {
			task.Date = next
		}
	}

	return nil
}

func (uc *TaskUseCase) afterNow(date, now time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}
