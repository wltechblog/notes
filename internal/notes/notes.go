package notes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wltechblog/notes/internal/platform"
)

type Note struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
}

func (n *Note) Path(baseDir string) string {
	return filepath.Join(baseDir, n.ID+".txt")
}

type NoteManager struct {
	baseDir string
}

func NewNoteManager() (*NoteManager, error) {
	baseDir, err := platform.GetDataDir(platform.NotesSubdir)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(baseDir, platform.GetDataDirPerm()); err != nil {
		return nil, fmt.Errorf("failed to create notes directory: %w", err)
	}

	return &NoteManager{baseDir: baseDir}, nil
}

func (nm *NoteManager) ListNotes() ([]Note, error) {
	var notes []Note

	entries, err := os.ReadDir(nm.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".txt")
		note, err := nm.loadNote(id)
		if err != nil {
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (nm *NoteManager) CreateNote(name string, content string) (*Note, error) {
	if name == "" {
		name = strconv.FormatInt(time.Now().Unix(), 10)
	}

	timestamp := time.Now()
	id, err := nm.getNextID()
	if err != nil {
		return nil, err
	}

	note := &Note{
		ID:        id,
		Name:      name,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Content:   content,
	}

	if err := nm.saveNote(note); err != nil {
		return nil, err
	}

	return note, nil
}

func (nm *NoteManager) GetNote(id string) (*Note, error) {
	note, err := nm.loadNote(id)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (nm *NoteManager) UpdateNote(id string, content string) (*Note, error) {
	note, err := nm.loadNote(id)
	if err != nil {
		return nil, err
	}

	note.Content = content
	note.UpdatedAt = time.Now()

	if err := nm.saveNote(&note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (nm *NoteManager) DeleteNote(id string) error {
	notePath := filepath.Join(nm.baseDir, id+".txt")
	if err := os.Remove(notePath); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}

func (nm *NoteManager) SearchNotes(keyword string) ([]Note, error) {
	notes, err := nm.ListNotes()
	if err != nil {
		return nil, err
	}

	var matchingNotes []Note
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Content), strings.ToLower(keyword)) ||
			strings.Contains(strings.ToLower(note.Name), strings.ToLower(keyword)) {
			matchingNotes = append(matchingNotes, note)
		}
	}

	return matchingNotes, nil
}

func (nm *NoteManager) loadNote(id string) (Note, error) {
	notePath := filepath.Join(nm.baseDir, id+".txt")
	data, err := os.ReadFile(notePath)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read note: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 {
		return Note{}, fmt.Errorf("invalid note format")
	}

	createdAt, err := time.Parse(time.RFC3339, strings.TrimPrefix(lines[0], "Created: "))
	if err != nil {
		return Note{}, fmt.Errorf("failed to parse created timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, strings.TrimPrefix(lines[1], "Updated: "))
	if err != nil {
		return Note{}, fmt.Errorf("failed to parse updated timestamp: %w", err)
	}

	name := strings.TrimPrefix(lines[2], "Name: ")
	content := strings.Join(lines[3:], "\n")

	return Note{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Content:   content,
	}, nil
}

func (nm *NoteManager) saveNote(note *Note) error {
	var content string
	content += fmt.Sprintf("Created: %s\n", note.CreatedAt.Format(time.RFC3339))
	content += fmt.Sprintf("Updated: %s\n", note.UpdatedAt.Format(time.RFC3339))
	content += fmt.Sprintf("Name: %s\n", note.Name)
	content += note.Content

	notePath := filepath.Join(nm.baseDir, note.ID+".txt")
	if err := os.WriteFile(notePath, []byte(content), platform.GetDataFilePerm()); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	return nil
}

func (nm *NoteManager) EditInEditor(note *Note) error {
	editor := platform.GetDefaultEditor()

	tmpFile, err := os.CreateTemp("", "note-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(note.Content); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	cmdArgs := platform.GetEditorArgs(editor, tmpPath)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" && platform.IsGUIEditor(editor) {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		cmd = exec.Command(editor, tmpPath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	editedContent, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to read edited content: %w", err)
	}

	note.Content = string(editedContent)
	note.UpdatedAt = time.Now()

	return nm.saveNote(note)
}

func (nm *NoteManager) getNextID() (string, error) {
	counterFile := filepath.Join(nm.baseDir, ".counter")

	data, err := os.ReadFile(counterFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to read counter file: %w", err)
		}
		data = []byte("0")
	}

	currentID, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse counter: %w", err)
	}

	nextID := currentID + 1

	if err := os.WriteFile(counterFile, []byte(strconv.FormatInt(nextID, 10)), platform.GetDataFilePerm()); err != nil {
		return "", fmt.Errorf("failed to write counter: %w", err)
	}

	return strconv.FormatInt(nextID, 10), nil
}
