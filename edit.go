package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit an existing note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noteMode {
			return fmt.Errorf("this command is only available for notes, use 'task edit' instead")
		}
		nm, err := notes.NewNoteManager()
		if err != nil {
			return err
		}

		id := args[0]
		note, err := nm.GetNote(id)
		if err != nil {
			fmt.Printf("Note not found: %s\n", id)
			return nil
		}

		if err := nm.EditInEditor(note); err != nil {
			return err
		}

		_, err = nm.UpdateNote(id, note.Content)
		if err != nil {
			return err
		}

		fmt.Printf("Note updated: %s\n", id)
		return nil
	},
}

func init() {
	if noteMode {
		rootCmd.AddCommand(editCmd)
	}
}
