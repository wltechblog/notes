package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskEditCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a task's note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		note, err := nm.GetNote(task.NoteID)
		if err != nil {
			fmt.Printf("Note not found: %s\n", task.NoteID)
			return nil
		}

		if err := nm.EditInEditor(note); err != nil {
			return err
		}

		_, err = nm.UpdateNote(task.NoteID, note.Content)
		if err != nil {
			return err
		}

		fmt.Printf("Task note updated: %s\n", id)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskEditCmd)
	}
}
