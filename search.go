package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wltechblog/notes/internal/notes"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search notes by keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noteMode {
			return fmt.Errorf("this command is only available for notes, use 'task search' instead")
		}
		nm, err := notes.NewNoteManager()
		if err != nil {
			return err
		}

		keyword := args[0]
		notesList, err := nm.SearchNotes(keyword)
		if err != nil {
			return err
		}

		if len(notesList) == 0 {
			fmt.Printf("No notes found matching '%s'\n", keyword)
			return nil
		}

		for _, note := range notesList {
			fmt.Printf("%s | %s | Created: %s | Updated: %s\n",
				note.ID,
				note.Name,
				note.CreatedAt.Format("2006-01-02 15:04:05"),
				note.UpdatedAt.Format("2006-01-02 15:04:05"))
		}

		return nil
	},
}

func init() {
	if noteMode {
		rootCmd.AddCommand(searchCmd)
	}
}
