package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search tasks by keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tm, err := tasks.NewTaskManager()
		if err != nil {
			return err
		}

		keyword := args[0]
		taskList, err := tm.SearchTasks(keyword)
		if err != nil {
			return err
		}

		if len(taskList) == 0 {
			fmt.Printf("No tasks found matching '%s'\n", keyword)
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
	if taskMode {
		rootCmd.AddCommand(taskSearchCmd)
	}
}
