package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noteMode {
			return fmt.Errorf("this command is only available for notes, use 'task delete' instead")
		}
		nm, err := notes.NewNoteManager()
		if err != nil {
			return err
		}

		id := args[0]
		if err := nm.DeleteNote(id); err != nil {
			fmt.Printf("Failed to delete note: %v\n", err)
			return nil
		}

		fmt.Printf("Note deleted: %s\n", id)
		return nil
	},
}

func init() {
	if noteMode {
		rootCmd.AddCommand(deleteCmd)
	}
}
