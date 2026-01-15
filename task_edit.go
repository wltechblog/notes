package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskEditCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a task's note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !taskMode {
			return fmt.Errorf("this command is only available for tasks, use 'note edit' instead")
		}
		tm, err := tasks.NewTaskManager()
		if err != nil {
			return err
		}

		id := args[0]
		task, err := tm.GetTask(id)
		if err != nil {
			fmt.Printf("Task not found: %s\n", id)
			return nil
		}

		if err := tm.EditInEditor(task); err != nil {
			return err
		}

		fmt.Printf("Task updated: %s\n", id)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskEditCmd)
	}
}
