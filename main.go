package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	noteMode bool
	taskMode bool
)

func getRootCommand() *cobra.Command {
	progName := filepath.Base(os.Args[0])

	if progName == "task" || progName == "tasks" {
		taskMode = true
		return &cobra.Command{
			Use:   "task",
			Short: "A CLI application for managing tasks",
		}
	} else {
		noteMode = true
		return &cobra.Command{
			Use:   "note",
			Short: "A CLI application for organizing note taking",
		}
	}
}

var rootCmd = getRootCommand()

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
