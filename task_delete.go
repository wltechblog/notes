package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskDeleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"del", "rm"},
	Short:   "Delete a task",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !taskMode {
			return fmt.Errorf("this command is only available for tasks, use 'note delete' instead")
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

		nm, err := notes.NewNoteManager()
		if err != nil {
			return err
		}

		if task.NoteID != "" {
			if err := nm.DeleteNote(task.NoteID); err != nil {
				fmt.Printf("Failed to delete note: %v\n", err)
			}
		}

		if err := tm.DeleteTask(id); err != nil {
			fmt.Printf("Failed to delete task: %v\n", err)
			return nil
		}

		fmt.Printf("Task deleted: %s\n", id)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskDeleteCmd)
	}
}
