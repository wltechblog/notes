package tasks

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

type Status string

const (
	StatusOpen      Status = "open"
	StatusCompleted Status = "completed"
	StatusAbandoned Status = "abandoned"
)

type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	NoteID    string    `json:"note_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
}

type TaskManager struct {
	baseDir string
}

func NewTaskManager() (*TaskManager, error) {
	baseDir, err := platform.GetDataDir(platform.TasksSubdir)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(baseDir, platform.GetDataDirPerm()); err != nil {
		return nil, fmt.Errorf("failed to create tasks directory: %w", err)
	}

	return &TaskManager{baseDir: baseDir}, nil
}

func (tm *TaskManager) ListTasks(statusFilter Status) ([]Task, error) {
	var tasks []Task

	entries, err := os.ReadDir(tm.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".txt")
		task, err := tm.loadTask(id)
		if err != nil {
			continue
		}

		if statusFilter == "" || task.Status == statusFilter {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (tm *TaskManager) CreateTask(name string, content string) (*Task, error) {
	if name == "" {
		name = strconv.FormatInt(time.Now().Unix(), 10)
	}

	timestamp := time.Now()
	id, err := tm.getNextID()
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID:        id,
		Name:      name,
		Status:    StatusOpen,
		NoteID:    "",
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
		Content:   content,
	}

	if err := tm.saveTask(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (tm *TaskManager) GetTask(id string) (*Task, error) {
	task, err := tm.loadTask(id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (tm *TaskManager) UpdateTask(id string, content string) (*Task, error) {
	task, err := tm.loadTask(id)
	if err != nil {
		return nil, err
	}

	task.Content = content
	task.UpdatedAt = time.Now()

	if err := tm.saveTask(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (tm *TaskManager) UpdateTaskStatus(id string, status Status) (*Task, error) {
	task, err := tm.loadTask(id)
	if err != nil {
		return nil, err
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	if err := tm.saveTask(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (tm *TaskManager) SearchTasks(keyword string) ([]Task, error) {
	tasks, err := tm.ListTasks("")
	if err != nil {
		return nil, err
	}

	var matchingTasks []Task
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Content), strings.ToLower(keyword)) ||
			strings.Contains(strings.ToLower(task.Name), strings.ToLower(keyword)) {
			matchingTasks = append(matchingTasks, task)
		}
	}

	return matchingTasks, nil
}

func (tm *TaskManager) loadTask(id string) (Task, error) {
	taskPath := filepath.Join(tm.baseDir, id+".txt")
	data, err := os.ReadFile(taskPath)
	if err != nil {
		return Task{}, fmt.Errorf("failed to read task: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 5 {
		return Task{}, fmt.Errorf("invalid task format")
	}

	createdAt, err := time.Parse(time.RFC3339, strings.TrimPrefix(lines[0], "Created: "))
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse created timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, strings.TrimPrefix(lines[1], "Updated: "))
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse updated timestamp: %w", err)
	}

	status := Status(strings.TrimPrefix(lines[2], "Status: "))
	noteID := strings.TrimPrefix(lines[3], "NoteID: ")
	name := strings.TrimPrefix(lines[4], "Name: ")
	content := strings.Join(lines[5:], "\n")

	return Task{
		ID:        id,
		Name:      name,
		Status:    status,
		NoteID:    noteID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Content:   content,
	}, nil
}

func (tm *TaskManager) saveTask(task *Task) error {
	var content string
	content += fmt.Sprintf("Created: %s\n", task.CreatedAt.Format(time.RFC3339))
	content += fmt.Sprintf("Updated: %s\n", task.UpdatedAt.Format(time.RFC3339))
	content += fmt.Sprintf("Status: %s\n", task.Status)
	content += fmt.Sprintf("NoteID: %s\n", task.NoteID)
	content += fmt.Sprintf("Name: %s\n", task.Name)
	content += task.Content

	taskPath := filepath.Join(tm.baseDir, task.ID+".txt")
	if err := os.WriteFile(taskPath, []byte(content), platform.GetDataFilePerm()); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

func (tm *TaskManager) getNextID() (string, error) {
	counterFile := filepath.Join(tm.baseDir, ".counter")

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

func (tm *TaskManager) EditInEditor(task *Task) error {
	editor := platform.GetDefaultEditor()

	tmpFile, err := os.CreateTemp("", "task-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(task.Content); err != nil {
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

	task.Content = string(editedContent)
	task.UpdatedAt = time.Now()

	_, err = tm.UpdateTask(task.ID, task.Content)
	return err
}

func (tm *TaskManager) DeleteTask(id string) error {
	task, err := tm.GetTask(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.NoteID != "" {
		notesDir, err := platform.GetDataDir(platform.NotesSubdir)
		if err != nil {
			return fmt.Errorf("failed to get notes directory: %w", err)
		}
		notePath := filepath.Join(notesDir, task.NoteID+".txt")
		if err := os.Remove(notePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete associated note: %w", err)
		}
	}

	taskPath := filepath.Join(tm.baseDir, id+".txt")
	if err := os.Remove(taskPath); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
