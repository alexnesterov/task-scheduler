// Package entity
package entity

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrCreateTask = errors.New("create task")
	ErrUpdateTask = errors.New("update task")
	ErrDoneTask   = errors.New("done task")
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
