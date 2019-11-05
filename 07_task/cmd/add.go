package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/maximepeschard/gophercises/07_task/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todoList := task.GetTodoList()
		if todoList == nil {
			return errors.New("TODO list not initialized")
		}

		t := task.NewTask(strings.Join(args, " "))
		err := todoList.AddTask(t)
		if err != nil {
			return err
		}

		fmt.Printf("Added \"%s\" to your TODO list.\n", t.Name)
		return nil
	},
}
