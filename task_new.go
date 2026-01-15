package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
	"github.com/wltechblog/notes/internal/tasks"
)

var taskNewCmd = &cobra.Command{
	Use:     "new [name]",
	Aliases: []string{"create"},
	Short:   "Create a new task",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		fmt.Printf("Task created: %s\n", task.ID)
		return nil
	},
}

func init() {
	if taskMode {
		rootCmd.AddCommand(taskNewCmd)
	}
}
