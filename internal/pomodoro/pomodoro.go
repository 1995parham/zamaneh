// Package pomodoro defines the timer domain: phases, durations, and sessions.
package pomodoro

import "time"

// Phase is one stage of the pomodoro cycle.
type Phase string

// Phases of one pomodoro cycle.
const (
	PhaseWork       Phase = "work"
	PhaseShortBreak Phase = "short_break"
	PhaseLongBreak  Phase = "long_break"
)

// Title returns a human-readable phase label.
func (p Phase) Title() string {
	switch p {
	case PhaseWork:
		return "Work"
	case PhaseShortBreak:
		return "Short Break"
	case PhaseLongBreak:
		return "Long Break"
	default:
		return string(p)
	}
}

// Config holds the durations and cadence of the pomodoro cycle.
type Config struct {
	Work           time.Duration
	ShortBreak     time.Duration
	LongBreak      time.Duration
	LongBreakEvery int
}

// DefaultConfig is the classic 25/5/15 pomodoro cadence.
func DefaultConfig() Config {
	return Config{
		Work:           25 * time.Minute,
		ShortBreak:     5 * time.Minute,
		LongBreak:      15 * time.Minute,
		LongBreakEvery: 4,
	}
}

// Duration returns how long the given phase should last under cfg.
func (c Config) Duration(p Phase) time.Duration {
	switch p {
	case PhaseWork:
		return c.Work
	case PhaseShortBreak:
		return c.ShortBreak
	case PhaseLongBreak:
		return c.LongBreak
	default:
		return c.Work
	}
}

// Next returns the phase that follows p, given how many work phases have been
// completed (including p, if p == PhaseWork and completed before incrementing).
func (c Config) Next(p Phase, workCount int) Phase {
	if p != PhaseWork {
		return PhaseWork
	}
	if c.LongBreakEvery > 0 && workCount%c.LongBreakEvery == 0 {
		return PhaseLongBreak
	}
	return PhaseShortBreak
}
