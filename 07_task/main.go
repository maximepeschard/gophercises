package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/maximepeschard/gophercises/07_task/cmd"
	"github.com/maximepeschard/gophercises/07_task/task"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	storageDir := filepath.Join(homeDir, ".local", "share", "task")
	err = os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(filepath.Join(storageDir, "task.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = task.InitTodoList(db, "task", "todo", "tasks")
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
