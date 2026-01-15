package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/tasks"
)

var (
	statusFilter string
)

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		tm, err := tasks.NewTaskManager()
		if err != nil {
			return err
		}

		var filter tasks.Status
		if statusFilter != "" {
			filter = tasks.Status(statusFilter)
			switch filter {
			case tasks.StatusOpen, tasks.StatusCompleted, tasks.StatusAbandoned:
			default:
				return fmt.Errorf("invalid status: %s (must be: open, completed, or abandoned)", statusFilter)
			}
		}

		taskList, err := tm.ListTasks(filter)
		if err != nil {
			return err
		}

		if len(taskList) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		for _, task := range taskList {
			fmt.Printf("%s | %s | [%s] | Note: %s | Created: %s | Updated: %s\n",
				task.ID,
				task.Name,
				task.Status,
				task.NoteID,
				task.CreatedAt.Format("2006-01-02 15:04:05"),
				task.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		return nil
	},
}

func init() {
	taskListCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "Filter by status (open, completed, abandoned)")
	if taskMode {
		rootCmd.AddCommand(taskListCmd)
	}
}
