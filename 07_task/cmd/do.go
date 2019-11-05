package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/maximepeschard/gophercises/07_task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskIndex, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		todoList := task.GetTodoList()
		if todoList == nil {
			return errors.New("TODO list not initialized")
		}

		t, err := todoList.MarkTaskComplete(taskIndex - 1)
		if err != nil {
			return err
		}

		fmt.Printf("You have completed the task \"%s\"\n", t.Name)
		return nil
	},
}
