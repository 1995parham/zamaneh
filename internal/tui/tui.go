// Package tui implements the Bubble Tea model for the zamaneh pomodoro timer.
package tui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/1995parham/zamaneh/internal/pomodoro"
	"github.com/1995parham/zamaneh/internal/store"
)

// tickMsg fires once per second while the timer is running.
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Model is the Bubble Tea model for the pomodoro timer.
type Model struct {
	cfg   pomodoro.Config
	title string

	phase     pomodoro.Phase
	remaining time.Duration
	startedAt time.Time

	paused    bool
	completed int // completed work phases
	width     int
	height    int

	progress progress.Model

	// title-edit mode
	editing bool
	input   textinput.Model

	dataPath string
	saveErr  error
	quitting bool
}

// New builds a model with a sensible initial state.
func New(title string, cfg pomodoro.Config, dataPath string) Model {
	p := progress.New(
		progress.WithDefaultBlend(),
		progress.WithWidth(48),
		progress.WithoutPercentage(),
	)

	ti := textinput.New()
	ti.Placeholder = "what are you working on?"
	ti.CharLimit = 80
	ti.SetWidth(40)
	ti.Prompt = "› "

	return Model{
		cfg:       cfg,
		title:     title,
		phase:     pomodoro.PhaseWork,
		remaining: cfg.Work,
		startedAt: time.Now(),
		progress:  p,
		input:     ti,
		dataPath:  dataPath,
	}
}

// Init starts the per-second timer.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), textinput.Blink)
}

// Update advances the state machine in response to messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := min(max(msg.Width-8, 20), 60)
		m.progress.SetWidth(w)
		return m, nil

	case tea.KeyPressMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				if v := strings.TrimSpace(m.input.Value()); v != "" {
					m.title = v
				}
				m.stopEditing()
				return m, nil
			case "esc":
				m.stopEditing()
				return m, nil
			case "ctrl+c":
				m.recordSession(false)
				m.quitting = true
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.recordSession(false)
			m.quitting = true
			return m, tea.Quit
		case "space", "p":
			m.paused = !m.paused
			return m, nil
		case "s":
			// skip current phase
			m.recordSession(false)
			m.advance()
			return m, nil
		case "r":
			// reset current phase to full duration
			m.remaining = m.cfg.Duration(m.phase)
			m.startedAt = time.Now()
			m.paused = false
			return m, nil
		case "t":
			m.startEditing()
			return m, textinput.Blink
		}
		return m, nil

	case tickMsg:
		if m.paused || m.editing || m.quitting {
			return m, tick()
		}
		m.remaining -= time.Second
		if m.remaining <= 0 {
			m.remaining = 0
			m.recordSession(true)
			m.advance()
		}
		return m, tick()

	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm
		return m, cmd
	}

	if m.editing {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) startEditing() {
	m.editing = true
	cur := m.title
	if cur == "-" {
		cur = ""
	}
	m.input.SetValue(cur)
	m.input.CursorEnd()
	m.input.Focus()
}

func (m *Model) stopEditing() {
	m.editing = false
	m.input.Blur()
	m.input.SetValue("")
}

// advance moves to the next phase in the cycle.
func (m *Model) advance() {
	if m.phase == pomodoro.PhaseWork {
		m.completed++
	}
	m.phase = m.cfg.Next(m.phase, m.completed)
	m.remaining = m.cfg.Duration(m.phase)
	m.startedAt = time.Now()
	m.paused = false
}

// recordSession appends the just-finished interval to the on-disk log.
func (m *Model) recordSession(completed bool) {
	dur := m.cfg.Duration(m.phase) - m.remaining
	if dur <= 0 {
		return
	}
	end := time.Now()
	s := store.Session{
		Title:     m.title,
		Phase:     m.phase,
		StartedAt: m.startedAt,
		EndedAt:   end,
		Duration:  dur,
		Completed: completed,
	}
	if _, err := store.Append(&s); err != nil {
		m.saveErr = err
	}
}

// ---- styling ---------------------------------------------------------------

var (
	colorPrimary = lipgloss.Color("#F5A8C1") // soft pink
	colorWork    = lipgloss.Color("#FF6B6B") // tomato red
	colorBreak   = lipgloss.Color("#4ECDC4") // teal
	colorLong    = lipgloss.Color("#7C7CF8") // periwinkle
	colorMuted   = lipgloss.Color("#6B7280")
	colorAccent  = lipgloss.Color("#FBBF24") // amber for paused
	colorText    = lipgloss.Color("#E5E7EB")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			Padding(0, 1)

	timerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			Padding(1, 4).
			MarginTop(1)

	pausedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	keyStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Bold(true)

	statStyle = lipgloss.NewStyle().Foreground(colorText)
)

func phaseColor(p pomodoro.Phase) color.Color {
	switch p {
	case pomodoro.PhaseWork:
		return colorWork
	case pomodoro.PhaseShortBreak:
		return colorBreak
	case pomodoro.PhaseLongBreak:
		return colorLong
	default:
		return colorText
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d / time.Second)
	mm := total / 60
	ss := total % 60
	return fmt.Sprintf("%02d:%02d", mm, ss)
}

// View renders the full frame.
func (m Model) View() tea.View {
	if m.quitting {
		v := tea.NewView("")
		return v
	}

	pc := phaseColor(m.phase)

	header := titleStyle.Render("✻ Zamaneh")
	if m.title != "" && m.title != "-" {
		header += subtitleStyle.Render("— " + m.title)
	}

	phaseTag := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0F172A")).
		Background(pc).
		Padding(0, 2).
		Render(strings.ToUpper(m.phase.Title()))

	state := lipgloss.NewStyle().Foreground(colorMuted).Render("running")
	switch {
	case m.editing:
		state = pausedStyle.Render("editing")
	case m.paused:
		state = pausedStyle.Render("paused")
	}

	statusRow := lipgloss.JoinHorizontal(lipgloss.Top, phaseTag, "  ", state)

	timeText := timerStyle.
		Foreground(pc).
		Render(formatDuration(m.remaining))

	total := m.cfg.Duration(m.phase)
	pct := 0.0
	if total > 0 {
		pct = 1.0 - float64(m.remaining)/float64(total)
		if pct < 0 {
			pct = 0
		}
		if pct > 1 {
			pct = 1
		}
	}
	bar := m.progress.ViewAs(pct)

	cycle := m.cfg.LongBreakEvery
	if cycle <= 0 {
		cycle = 4
	}
	currentInCycle := m.completed%cycle + 1
	if m.phase != pomodoro.PhaseWork {
		currentInCycle = m.completed % cycle
		if currentInCycle == 0 {
			currentInCycle = cycle
		}
	}
	pips := renderPips(cycle, currentInCycle, pc)

	stats := statStyle.Render(fmt.Sprintf(
		"completed pomodoros: %d   cycle: %d/%d",
		m.completed, currentInCycle, cycle,
	))

	hints := footerStyle.Render(
		keyStyle.Render("space") + " pause  " +
			keyStyle.Render("s") + " skip  " +
			keyStyle.Render("r") + " reset  " +
			keyStyle.Render("t") + " title  " +
			keyStyle.Render("q") + " quit",
	)
	if m.editing {
		hints = footerStyle.Render(
			keyStyle.Render("enter") + " save  " +
				keyStyle.Render("esc") + " cancel",
		)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pc).
		Padding(1, 3).
		Align(lipgloss.Center)

	rows := []string{
		statusRow,
		timeText,
		bar,
		"",
		pips,
		stats,
	}
	if m.editing {
		editStyle := lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)
		rows = append(rows,
			editStyle.Render("rename session:"),
			m.input.View(),
		)
	}
	rows = append(rows, hints)
	body := lipgloss.JoinVertical(lipgloss.Center, rows...)
	content := lipgloss.JoinVertical(lipgloss.Center, header, box.Render(body))

	if m.saveErr != nil {
		content = lipgloss.JoinVertical(lipgloss.Center,
			content,
			errStyle.Render("save error: "+m.saveErr.Error()),
		)
	}

	// Center on screen if we know the size.
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center, content)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "Zamaneh — " + m.phase.Title()
	return v
}

func renderPips(total, current int, accent color.Color) string {
	on := lipgloss.NewStyle().Foreground(accent).Render("●")
	off := lipgloss.NewStyle().Foreground(colorMuted).Render("○")
	parts := make([]string, total)
	for i := range total {
		if i < current {
			parts[i] = on
		} else {
			parts[i] = off
		}
	}
	return strings.Join(parts, " ")
}

// Farewell is shown after the program exits — used by main to print a summary
// to the normal terminal after leaving the alt screen.
func (m Model) Farewell() string {
	return fmt.Sprintf(
		"you worked on %q · completed %d pomodoro(s) · log: %s",
		m.title, m.completed, m.dataPath,
	)
}
