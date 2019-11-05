package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/maximepeschard/gophercises/07_task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task from your TODO list",
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

		t, err := todoList.RemoveTask(taskIndex - 1)
		if err != nil {
			return err
		}

		fmt.Printf("You have deleted the task \"%s\"\n", t.Name)
		return nil
	},
}
