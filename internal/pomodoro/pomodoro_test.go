package pomodoro_test

import (
	"testing"
	"time"

	"github.com/1995parham/zamaneh/internal/pomodoro"
)

func TestDefaultConfigDurations(t *testing.T) {
	c := pomodoro.DefaultConfig()
	if c.Work != 25*time.Minute {
		t.Errorf("work: got %v, want 25m", c.Work)
	}
	if c.ShortBreak != 5*time.Minute {
		t.Errorf("short break: got %v, want 5m", c.ShortBreak)
	}
	if c.LongBreak != 15*time.Minute {
		t.Errorf("long break: got %v, want 15m", c.LongBreak)
	}
	if c.LongBreakEvery != 4 {
		t.Errorf("long break every: got %d, want 4", c.LongBreakEvery)
	}
}

func TestNext(t *testing.T) {
	c := pomodoro.DefaultConfig()

	tests := []struct {
		name      string
		phase     pomodoro.Phase
		workCount int
		want      pomodoro.Phase
	}{
		{"after 1st work → short", pomodoro.PhaseWork, 1, pomodoro.PhaseShortBreak},
		{"after 2nd work → short", pomodoro.PhaseWork, 2, pomodoro.PhaseShortBreak},
		{"after 3rd work → short", pomodoro.PhaseWork, 3, pomodoro.PhaseShortBreak},
		{"after 4th work → long", pomodoro.PhaseWork, 4, pomodoro.PhaseLongBreak},
		{"after 8th work → long", pomodoro.PhaseWork, 8, pomodoro.PhaseLongBreak},
		{"after short break → work", pomodoro.PhaseShortBreak, 1, pomodoro.PhaseWork},
		{"after long break → work", pomodoro.PhaseLongBreak, 4, pomodoro.PhaseWork},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.Next(tc.phase, tc.workCount); got != tc.want {
				t.Errorf("Next(%s, %d) = %s, want %s", tc.phase, tc.workCount, got, tc.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	c := pomodoro.DefaultConfig()
	if c.Duration(pomodoro.PhaseWork) != c.Work {
		t.Errorf("Duration(work) mismatch")
	}
	if c.Duration(pomodoro.PhaseShortBreak) != c.ShortBreak {
		t.Errorf("Duration(short) mismatch")
	}
	if c.Duration(pomodoro.PhaseLongBreak) != c.LongBreak {
		t.Errorf("Duration(long) mismatch")
	}
}
