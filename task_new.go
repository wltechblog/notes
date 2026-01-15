package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskNewCmd = &cobra.Command{
	Use:     "new [name]",
	Aliases: []string{"create"},
	Short:   "Create a new task",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !taskMode {
			return fmt.Errorf("this command is only available for tasks, use 'note new' instead")
		}
		tm, err := tasks.NewTaskManager()
		if err != nil {
			return err
		}

		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		task, err := tm.CreateTask(name, "")
		if err != nil {
			return err
		}

		if err := tm.EditInEditor(task); err != nil {
			return err
		}

		if task.Content == "" {
			if err := tm.DeleteTask(task.ID); err != nil {
				return err
			}
			fmt.Println("Task not saved (empty content)")
			return nil
		}

		fmt.Printf("Task created: %s\n", task.ID)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskNewCmd)
	}
}
