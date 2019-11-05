package task

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// TodoList provides methods to manage tasks of a TODO list
type TodoList interface {
	// AddTask adds a new task to the TODO list
	AddTask(t *Task) error
	// ListTasks returns a list of all the imcomplete tasks in the TODO list
	ListTasks() ([]*Task, error)
	// ListTasksFilter returns a list of all the tasks in the TODO list verifying a predicate
	ListTasksFilter(fn func(*Task) bool) ([]*Task, error)
	// MarkTaskComplete marks the task in index position in the TODO list as complete
	MarkTaskComplete(index int) (*Task, error)
	// RemoveTask removes the task in index position from the TODO list
	RemoveTask(index int) (*Task, error)
}

type todoList struct {
	db          *bolt.DB
	rootBucket  []byte
	todoKey     []byte
	tasksBucket []byte
}

func (tl *todoList) getTask(tx *bolt.Tx, taskUID string) (*Task, error) {
	var t Task
	rootBucket := tx.Bucket(tl.rootBucket)
	tasksBucket := rootBucket.Bucket(tl.tasksBucket)

	taskPayload := tasksBucket.Get([]byte(taskUID))
	if taskPayload == nil {
		return nil, nil
	}

	err := json.Unmarshal(taskPayload, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (tl *todoList) putTask(tx *bolt.Tx, t *Task) error {
	rootBucket := tx.Bucket(tl.rootBucket)
	tasksBucket := rootBucket.Bucket(tl.tasksBucket)

	enc, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return tasksBucket.Put([]byte(t.UID), enc)
}

func (tl *todoList) deleteTask(tx *bolt.Tx, taskUID string) error {
	rootBucket := tx.Bucket(tl.rootBucket)
	tasksBucket := rootBucket.Bucket(tl.tasksBucket)

	return tasksBucket.Delete([]byte(taskUID))
}

func (tl *todoList) getTodo(tx *bolt.Tx) ([]string, error) {
	var todo []string
	rootBucket := tx.Bucket(tl.rootBucket)

	todoBytes := rootBucket.Get([]byte(tl.todoKey))
	if todoBytes == nil {
		return nil, nil
	}

	err := json.Unmarshal(todoBytes, &todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (tl *todoList) putTodo(tx *bolt.Tx, todo []string) error {
	rootBucket := tx.Bucket(tl.rootBucket)

	enc, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return rootBucket.Put(tl.todoKey, enc)
}

func (tl *todoList) AddTask(t *Task) error {
	return tl.db.Update(func(tx *bolt.Tx) error {
		err := tl.putTask(tx, t)
		if err != nil {
			return err
		}

		todo, err := tl.getTodo(tx)
		if err != nil {
			return err
		}

		todo = append(todo, t.UID)
		return tl.putTodo(tx, todo)
	})
}

func (tl *todoList) ListTasks() ([]*Task, error) {
	var tasks []*Task

	err := tl.db.View(func(tx *bolt.Tx) error {
		todo, err := tl.getTodo(tx)
		if err != nil {
			return err
		}

		for _, taskUID := range todo {
			task, err := tl.getTask(tx, taskUID)
			if err != nil {
				return err
			}

			tasks = append(tasks, task)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tl *todoList) ListTasksFilter(fn func(*Task) bool) ([]*Task, error) {
	var tasks []*Task

	err := tl.db.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket(tl.rootBucket)
		tasksBucket := rootBucket.Bucket(tl.tasksBucket)

		tasksBucket.ForEach(func(k, v []byte) error {
			var t Task
			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}

			if fn(&t) {
				tasks = append(tasks, &t)
			}

			return nil
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tl *todoList) MarkTaskComplete(index int) (*Task, error) {
	var t *Task

	err := tl.db.Update(func(tx *bolt.Tx) error {
		todo, err := tl.getTodo(tx)
		if err != nil {
			return err
		}

		taskUID := todo[index]
		todo = removeElement(todo, index)
		err = tl.putTodo(tx, todo)
		if err != nil {
			return err
		}

		t, err = tl.getTask(tx, taskUID)
		if err != nil {
			return err
		}

		t.Complete = true
		t.CompletionTime = time.Now().Unix()
		return tl.putTask(tx, t)
	})

	return t, err
}

func (tl *todoList) RemoveTask(index int) (*Task, error) {
	var t *Task

	err := tl.db.Update(func(tx *bolt.Tx) error {
		todo, err := tl.getTodo(tx)
		if err != nil {
			return err
		}

		taskUID := todo[index]
		todo = removeElement(todo, index)
		err = tl.putTodo(tx, todo)
		if err != nil {
			return err
		}

		t, err = tl.getTask(tx, taskUID)

		return tl.deleteTask(tx, taskUID)
	})

	return t, err
}

func removeElement(l []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, l[:index]...)
	return append(ret, l[index+1:]...)
}

// The TODO list is a package level variable
var tl TodoList

// InitTodoList initializes the TODO list given a database and keys
func InitTodoList(db *bolt.DB, bucket, todoKey, tasksBucket string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		_, err = root.CreateBucketIfNotExists([]byte(tasksBucket))
		return err
	})
	if err != nil {
		return err
	}
	tl = &todoList{
		db:          db,
		rootBucket:  []byte(bucket),
		todoKey:     []byte(todoKey),
		tasksBucket: []byte(tasksBucket),
	}

	return nil
}

// GetTodoList returns the TODO list
func GetTodoList() TodoList {
	return tl
}
