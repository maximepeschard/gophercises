package task

import (
	"github.com/google/uuid"
)

// Task represents a n item of a TODO list
type Task struct {
	Name           string
	Complete       bool
	CompletionTime int64
	UID            string
}

// NewTask returns a task with a name and an UID
func NewTask(name string) *Task {
	return &Task{Name: name, UID: uuid.New().String()}
}
