package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wltechblog/notes/internal/notes"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".local", "share", "tasks")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
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

	nm, err := notes.NewNoteManager()
	if err != nil {
		return nil, err
	}

	note, err := nm.CreateNote("Task: "+name, content)
	if err != nil {
		return nil, err
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
		NoteID:    note.ID,
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

func (tm *TaskManager) DeleteTask(id string) error {
	taskPath := filepath.Join(tm.baseDir, id+".txt")
	if err := os.Remove(taskPath); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
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
	if len(lines) < 4 {
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

	var noteID string
	var name string
	var content string

	if len(lines) >= 5 && strings.HasPrefix(lines[3], "NoteID:") {
		noteID = strings.TrimPrefix(lines[3], "NoteID: ")
		name = strings.TrimPrefix(lines[4], "Name: ")
		content = strings.Join(lines[5:], "\n")
	} else {
		name = strings.TrimPrefix(lines[3], "Name: ")
		content = strings.Join(lines[4:], "\n")
	}

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
	if err := os.WriteFile(taskPath, []byte(content), 0644); err != nil {
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

	if err := os.WriteFile(counterFile, []byte(strconv.FormatInt(nextID, 10)), 0644); err != nil {
		return "", fmt.Errorf("failed to write counter: %w", err)
	}

	return strconv.FormatInt(nextID, 10), nil
}
