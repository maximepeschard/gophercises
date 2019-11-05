package cmd

import (
	"errors"
	"fmt"

	"github.com/maximepeschard/gophercises/07_task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoList := task.GetTodoList()
		if todoList == nil {
			return errors.New("TODO list not initialized")
		}

		tasks, err := todoList.ListTasks()
		if err != nil {
			return err
		}

		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i+1, t.Name)
		}

		return nil
	},
}
