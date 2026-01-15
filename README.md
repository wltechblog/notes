# notes

A command-line application for organizing and managing personal notes and tasks. Notes and tasks are stored as plain text files in the filesystem, each with a unique ID, name, created timestamp, and last edited timestamp. Tasks also include status tracking (open, completed, abandoned).

The same binary serves both functions - invoked as `note` for notes and `task` for tasks (via symlink).

## Features

- **Simple file-based storage**: Notes stored in `~/.local/share/notes/`, tasks in `~/.local/share/tasks/`
- **Unique identification**: Each note/task has a unique sequential numeric ID
- **Timestamp tracking**: Track when notes/tasks were created and last edited
- **Text search**: Search across note/task names and content
- **Editor integration**: Uses your `$EDITOR` for creating and editing notes
- **Task status management**: Change task status between open, completed, and abandoned
- **Status filtering**: List tasks by status
- **Shell completion**: Auto-complete support for bash, zsh, fish, and powershell

## Installation

### From source

```bash
git clone https://github.com/wltechblog/notes.git
cd notes
go build -o note
sudo cp note /usr/local/bin/  # or copy to ~/.local/bin/
sudo ln -s /usr/local/bin/note /usr/local/bin/task
```

### Using Makefile

```bash
make build    # Build binary
make install  # Install to ~/.local/bin/ (also creates 'task' symlink)
```

The installation creates both `note` and `task` commands. The `task` command is a symlink to the same binary.

## Task Commands

### Create a new task

```bash
task new "Buy groceries"    # Create task with name
task new                   # Create task with timestamp as name
```

Tasks start with `open` status by default.

### List all tasks

```bash
task list    # List all tasks
```

Outputs tasks in format: `ID | Name | [status] | Created: date | Updated: date`

### List tasks by status

```bash
task list --status open        # List only open tasks
task list --status completed   # List only completed tasks
task list --status abandoned  # List only abandoned tasks
task list -s completed       # Short form
```

### Change task status

```bash
task status <id> <status>    # Change task status
task status 1 completed      # Example: mark task 1 as completed
task status 2 abandoned      # Example: mark task 2 as abandoned
```

Valid statuses: `open`, `completed`, `abandoned`

### Search tasks

```bash
task search "keyword"    # Search by keyword
task search "meeting"    # Example: find all meeting tasks
```

Performs case-insensitive search across task names and content.

### Delete a task

```bash
task delete <id>         # Delete a task by ID
task delete 1            # Example: delete task 1
```

Deleting a task also deletes its associated note.

### Edit a task's note

```bash
task edit <id>          # Edit the note associated with a task
task edit 1            # Example: edit note for task 1
```

Opens `$EDITOR` with the note content associated with the task. Updates the note's content and last edited timestamp.

## Task & Note Integration

Tasks automatically create associated notes:

- When you create a task with `task new`, a corresponding note is created with name "Task: [task name]"
- The task's note ID is displayed in task listings
- Use `task edit <id>` to edit the note content of a task
- The task's note is visible in `note list` and searchable with `note search`
- Deleting a task also deletes its associated note

This keeps your task documentation separate from task status tracking.

## Note Commands

### Create a new note

```bash
note new "my note"     # Create note with custom name
note new               # Create note with timestamp as name
```

Opens `$EDITOR` (defaults to `vi` if not set) with an empty buffer. If you save with content, the note is created. If you exit with an empty buffer, no note is saved.

### List all notes

```bash
note list
```

Outputs notes in format: `ID | Name | Created: date | Updated: date`

### Search notes

```bash
note search "keyword"    # Search by keyword
note search "meeting"    # Example: find all meeting notes
```

Performs case-insensitive search across note names and content.

### Edit a note

```bash
note edit <id>          # Opens $EDITOR with note content
note edit a1b2c3d4      # Example: edit specific note
```

Updates the note's content and last edited timestamp when saved.

### Delete a note

```bash
note delete <id>         # Delete a note by ID
note delete a1b2c3d4    # Example: delete specific note
```

## Shell Completion

Enable command-line completion for your shell. The `task` command completion works the same way as `note`:

### Bash

Add to your `~/.bashrc` or `~/.bash_profile`:

```bash
# For auto-completion
source <(note completion bash)
source <(task completion bash)
```

Or for persistent completion:

```bash
note completion bash > ~/.local/share/bash-completion/completions/note
task completion bash > ~/.local/share/bash-completion/completions/task
```

### Zsh

Add to your `~/.zshrc`:

```bash
# For auto-completion
source <(note completion zsh)
source <(task completion zsh)

# Or add to your completion functions directory
note completion zsh > ~/.zsh/completion/_note
task completion zsh > ~/.zsh/completion/_task
```

### Fish

Add to your `~/.config/fish/completions/` directory:

```bash
note completion fish > ~/.config/fish/completions/note.fish
task completion fish > ~/.config/fish/completions/task.fish
```

### PowerShell

Add to your PowerShell profile:

```powershell
note completion powershell | Out-String | Invoke-Expression
task completion powershell | Out-String | Invoke-Expression
```

Or save and source from your profile:

```powershell
note completion powershell > note.ps1
task completion powershell > task.ps1
```

## Storage

### Note Storage

Notes are stored as plain text files in `~/.local/share/notes/`:

```
~/.local/share/notes/
├── 1.txt
├── 2.txt
├── 3.txt
├── .counter    # Tracks next ID
└── ...
```

Each note file contains:

```
Created: 2026-01-15T10:49:30-07:00
Updated: 2026-01-15T10:49:56-07:00
Name: my note
This is note content...
```

### Task Storage

Tasks are stored as plain text files in `~/.local/share/tasks/`:

```
~/.local/share/tasks/
├── 1.txt
├── 2.txt
├── 3.txt
├── .counter    # Tracks next ID
└── ...
```

Each task file contains:

```
Created: 2026-01-15T10:49:30-07:00
Updated: 2026-01-15T10:49:56-07:00
Status: open
NoteID: 7
Name: Buy groceries
This is task content...
```

The `NoteID` field references the associated note in `~/.local/share/notes/`.
```

Each note file contains:

```
Created: 2026-01-15T10:49:30-07:00
Updated: 2026-01-15T10:49:56-07:00
Name: my note
This is the note content...
```

## Configuration

- **Editor**: Set via `$EDITOR` environment variable (defaults to `vi`)
- **Storage location**: `~/.local/share/notes/` (XDG Base Directory compliant)
- **ID generation**: Sequential numbers stored in `.counter` file

## Development

### Project structure

```
.
├── main.go                      # Application entry point (detects note/task mode)
├── new.go                       # Create new notes command
├── list.go                      # List notes command
├── search.go                    # Search notes command
├── edit.go                      # Edit notes command
├── delete.go                    # Delete notes command
├── task_new.go                  # Create new tasks command
├── task_list.go                 # List tasks command
├── task_search.go               # Search tasks command
├── task_edit.go                # Edit task's note command
├── task_delete.go              # Delete tasks command
├── task_status.go              # Change task status command
├── internal/
│   ├── notes/
│   │   └── notes.go          # Core note logic and storage
│   └── tasks/
│       └── tasks.go          # Core task logic and storage
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Building

```bash
go build -o note      # Build binary
go test ./...         # Run tests
```

## License

This software is licensed under the GNU GPL 2.0.
