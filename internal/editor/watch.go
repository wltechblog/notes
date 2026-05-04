package editor

import (
	"context"
	"os"
	"time"
)

// Watch polls path for mtime changes and calls onSave with the file
// contents whenever the mtime advances. It returns when ctx is done.
// Errors reading the file are silently skipped so that editor UIs are
// not corrupted by terminal output during an active session.
func Watch(ctx context.Context, path string, interval time.Duration, onSave func(data []byte)) {
	var lastMtime time.Time
	if st, err := os.Stat(path); err == nil {
		lastMtime = st.ModTime()
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			st, err := os.Stat(path)
			if err != nil {
				continue
			}
			if !st.ModTime().After(lastMtime) {
				continue
			}
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			lastMtime = st.ModTime()
			onSave(data)
		}
	}
}
