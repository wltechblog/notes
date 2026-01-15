package platform

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	NotesSubdir = "notes"
	TasksSubdir = "tasks"
)

func GetDataDir(subdir string) (string, error) {
	var baseDir string

	if runtime.GOOS == "windows" {
		baseDir = os.Getenv("LOCALAPPDATA")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, ".local", "share")
	}

	if subdir != "" {
		baseDir = filepath.Join(baseDir, subdir)
	}

	return baseDir, nil
}

func GetDataDirPerm() os.FileMode {
	if runtime.GOOS == "windows" {
		return 0755
	}
	return 0755
}

func GetDataFilePerm() os.FileMode {
	if runtime.GOOS == "windows" {
		return 0644
	}
	return 0644
}

func GetDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vi"
		}
	}
	return editor
}

func GetEditorArgs(editor string, filepath string) []string {
	if runtime.GOOS == "windows" {
		switch strings.ToLower(editor) {
		case "notepad":
			return []string{"notepad", filepath}
		case "notepad++", "notepadplusplus":
			return []string{"notepad++", filepath}
		case "code", "vscode":
			return []string{"code", "--wait", filepath}
		case "subl", "sublime":
			return []string{"subl", "--wait", filepath}
		default:
			return []string{editor, filepath}
		}
	}
	return []string{editor, filepath}
}

func IsGUIEditor(editor string) bool {
	guiEditors := map[string]bool{
		"notepad":         true,
		"notepad++":       true,
		"notepadplusplus": true,
		"code":            true,
		"vscode":          true,
		"subl":            true,
		"sublime":         true,
		"vim":             false,
		"vi":              false,
		"nano":            false,
		"emacs":           false,
	}
	return guiEditors[strings.ToLower(editor)]
}
