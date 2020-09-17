package timer_test

import (
	"testing"
	"time"

	"github.com/1995parham/zamaneh/timer"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	tt := timer.New()

	d1 := <-tt.C
	t1 := time.Now()
	assert.Equal(t, d1, time.Duration(0))

	d2 := <-tt.C
	t2 := time.Now()
	assert.Equal(t, d2, time.Second)

	assert.True(t, t2.Sub(t1) <= time.Second+500*time.Microsecond)
}

func TestStop(t *testing.T) {
	tt := timer.New()

	d1 := <-tt.C
	assert.Equal(t, d1, time.Duration(0))

	d2 := <-tt.C
	tt.Stop()
	assert.Equal(t, d2, time.Second)

	time.Sleep(2 * time.Second)

	tt.Start()

	d3 := <-tt.C
	assert.Equal(t, d3, 2*time.Second)
}
