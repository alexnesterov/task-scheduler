// Package usecase
package usecase

import (
	"time"

	"github.com/alexnesterov/task-scheduler/internal/domain/entity"
	"github.com/alexnesterov/task-scheduler/internal/domain/port"
)

type TaskService struct {
	TaskRepo port.TaskRepository
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
