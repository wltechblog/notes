package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
)

var newCmd = &cobra.Command{
	Use:     "new [name]",
	Aliases: []string{"create"},
	Short:   "Create a new note",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noteMode {
			return fmt.Errorf("this command is only available for notes, use 'task new' instead")
		}
		nm, err := notes.NewNoteManager()
		if err != nil {
			return err
		}

		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		note, err := nm.CreateNote(name, "")
		if err != nil {
			return err
		}

		if err := nm.EditInEditor(note); err != nil {
			return err
		}

		if note.Content == "" {
			if err := nm.DeleteNote(note.ID); err != nil {
				return err
			}
			fmt.Println("Note not saved (empty content)")
			return nil
		}

		_, err = nm.UpdateNote(note.ID, note.Content)
		if err != nil {
			return err
		}

		fmt.Printf("Note created: %s\n", note.ID)
		return nil
	},
}

func init() {
	if noteMode {
		rootCmd.AddCommand(newCmd)
	}
}
