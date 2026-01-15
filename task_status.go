package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskStatusCmd = &cobra.Command{
	Use:   "status [id] [status]",
	Short: "Change task status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tm, err := tasks.NewTaskManager()
		if err != nil {
			return err
		}

		id := args[0]
		status := tasks.Status(args[1])

		switch status {
		case tasks.StatusOpen, tasks.StatusCompleted, tasks.StatusAbandoned:
		default:
			return fmt.Errorf("invalid status: %s (must be: open, completed, or abandoned)", status)
		}

		_, err = tm.UpdateTaskStatus(id, status)
		if err != nil {
			fmt.Printf("Task not found: %s\n", id)
			return nil
		}

		fmt.Printf("Task %s status updated to: %s\n", id, status)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskStatusCmd)
	}
}
