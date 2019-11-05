package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/maximepeschard/gophercises/07_task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completedCmd)
}

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all tasks completed in the last 24 hours",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoList := task.GetTodoList()
		if todoList == nil {
			return errors.New("TODO list not initialized")
		}

		tasks, err := todoList.ListTasksFilter(func(t *task.Task) bool {
			return t.Complete && time.Since(time.Unix(t.CompletionTime, 0)).Hours() <= 24
		})
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return nil
		}

		fmt.Println("Completed tasks in the last 24 hours :")
		for _, t := range tasks {
			fmt.Printf("✓ %s\n", t.Name)
		}

		return nil
	},
}
