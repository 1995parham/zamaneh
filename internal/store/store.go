// Package store persists pomodoro sessions to a TOML file under XDG_DATA_HOME.
package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"

	"github.com/1995parham/zamaneh/internal/pomodoro"
)

const (
	appDir   = "zamaneh"
	dataFile = "sessions.toml"
)

// Session is one completed (or cancelled) interval of the pomodoro cycle.
type Session struct {
	Title     string         `toml:"title"`
	Phase     pomodoro.Phase `toml:"phase"`
	StartedAt time.Time      `toml:"started_at"`
	EndedAt   time.Time      `toml:"ended_at"`
	Duration  time.Duration  `toml:"duration"`
	Completed bool           `toml:"completed"`
}

// Document is the on-disk schema. A flat `sessions` array of tables keeps the
// file diff-friendly and easy to read by hand.
type Document struct {
	Sessions []Session `toml:"sessions"`
}

// Path returns the absolute location of the sessions file, creating the parent
// directory if needed.
func Path() (string, error) {
	dir := filepath.Join(xdg.DataHome, appDir)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("create data dir: %w", err)
	}
	return filepath.Join(dir, dataFile), nil
}

// Load reads the sessions file. A missing file yields an empty document.
func Load() (Document, string, error) {
	path, err := Path()
	if err != nil {
		return Document{}, "", err
	}

	// path is derived from xdg.DataHome, not user input — G304 is a false positive here.
	data, err := os.ReadFile(path) //nolint:gosec // see comment above
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Document{}, path, nil
		}
		return Document{}, path, fmt.Errorf("read sessions: %w", err)
	}

	var doc Document
	if err := toml.Unmarshal(data, &doc); err != nil {
		return Document{}, path, fmt.Errorf("parse sessions: %w", err)
	}
	return doc, path, nil
}

// Append adds a session to the file and writes it atomically.
func Append(s *Session) (string, error) {
	doc, path, err := Load()
	if err != nil {
		return path, err
	}
	doc.Sessions = append(doc.Sessions, *s)

	out, err := toml.Marshal(doc)
	if err != nil {
		return path, fmt.Errorf("encode sessions: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".sessions-*.toml")
	if err != nil {
		return path, fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }

	if _, err := tmp.Write(out); err != nil {
		_ = tmp.Close()
		cleanup()
		return path, fmt.Errorf("write temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return path, fmt.Errorf("close temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return path, fmt.Errorf("rename temp: %w", err)
	}
	return path, nil
}
