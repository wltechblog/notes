package gitbackup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func run(dir string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func EnsureRepo(dir string) error {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return ensureIdentity(dir)
	}
	if out, err := run(dir, "init", "-q"); err != nil {
		return fmt.Errorf("git init: %w: %s", err, firstLine(out))
	}
	if err := ensureIdentity(dir); err != nil {
		return err
	}
	return Commit(dir, "initial state")
}

// ensureIdentity sets a local fallback identity only if none is configured anywhere.
func ensureIdentity(dir string) error {
	if out, err := run(dir, "config", "user.email"); err != nil {
		return fmt.Errorf("git config user.email: %w: %s", err, firstLine(out))
	} else if len(strings.TrimSpace(string(out))) == 0 {
		if out, err := run(dir, "config", "user.email", "notes@localhost"); err != nil {
			return fmt.Errorf("git config user.email: %w: %s", err, firstLine(out))
		}
	}
	if out, err := run(dir, "config", "user.name"); err != nil {
		return fmt.Errorf("git config user.name: %w: %s", err, firstLine(out))
	} else if len(strings.TrimSpace(string(out))) == 0 {
		if out, err := run(dir, "config", "user.name", "notes"); err != nil {
			return fmt.Errorf("git config user.name: %w: %s", err, firstLine(out))
		}
	}
	return nil
}

func Commit(dir, message string) error {
	if out, err := run(dir, "add", "-A"); err != nil {
		return fmt.Errorf("git add: %w: %s", err, firstLine(out))
	}
	if _, err := run(dir, "diff", "--cached", "--quiet"); err == nil {
		return nil
	}
	if out, err := run(dir, "commit", "--quiet", "-m", message); err != nil {
		return fmt.Errorf("git commit: %w: %s", err, firstLine(out))
	}
	return nil
}

func firstLine(b []byte) string {
	s := strings.TrimSpace(string(b))
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func Warn(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "backup warning: %v\n", err)
	}
}
