package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
