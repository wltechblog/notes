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
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all tasks",
	Long:    "List all tasks. Use --status flag to filter by open, completed, or abandoned",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !taskMode {
			return fmt.Errorf("this command is only available for tasks, use 'note list' instead")
		}
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
			contentPreview := task.Content
			if len(contentPreview) > 30 {
				contentPreview = contentPreview[:30] + "..."
			}
			fmt.Printf("%s | %s | [%s] | %s | Created: %s | Updated: %s\n",
				task.ID,
				task.Name,
				task.Status,
				contentPreview,
				task.CreatedAt.Format("2006-01-02 15:04:05"),
				task.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		return nil
	},
}

func init() {
	if taskMode {
		taskListCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "Filter by status (open, completed, abandoned)")
		rootCmd.AddCommand(taskListCmd)
	}
}
