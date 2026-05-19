// Zamaneh is a terminal pomodoro timer with a Bubble Tea TUI.
// Sessions are appended to a TOML file under $XDG_DATA_HOME/zamaneh/sessions.toml.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/1995parham/zamaneh/internal/pomodoro"
	"github.com/1995parham/zamaneh/internal/store"
	"github.com/1995parham/zamaneh/internal/tui"
)

const usage = `zamaneh — the way you manage your time

Tracks pomodoro cycles in your terminal. Work and break intervals alternate
automatically. Completed intervals are appended to a TOML log under
$XDG_DATA_HOME/zamaneh/sessions.toml.

usage:
  zamaneh [flags]

flags:
`

// Set via -ldflags at release time by GoReleaser.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg := pomodoro.DefaultConfig()

	title := flag.String("title", "-", "label for this working session")
	work := flag.Duration("work", cfg.Work, "duration of a work interval")
	short := flag.Duration("short-break", cfg.ShortBreak, "duration of a short break")
	long := flag.Duration("long-break", cfg.LongBreak, "duration of a long break")
	every := flag.Int("long-break-every", cfg.LongBreakEvery, "long break after this many work intervals")
	verbose := flag.Bool("verbose", false, "enable debug logging")
	showVersion := flag.Bool("version", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("zamaneh %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	logLevel := slog.LevelWarn
	if *verbose {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))

	cfg.Work = sanitize(*work, cfg.Work)
	cfg.ShortBreak = sanitize(*short, cfg.ShortBreak)
	cfg.LongBreak = sanitize(*long, cfg.LongBreak)
	if *every > 0 {
		cfg.LongBreakEvery = *every
	}

	path, err := store.Path()
	if err != nil {
		fmt.Fprintf(os.Stderr, "zamaneh: %v\n", err)
		os.Exit(1)
	}

	model := tui.New(*title, cfg, path)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "zamaneh: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(tui.Model); ok {
		fmt.Println(m.Farewell())
	}
}

func sanitize(d, fallback time.Duration) time.Duration {
	if d <= 0 {
		return fallback
	}
	return d
}
