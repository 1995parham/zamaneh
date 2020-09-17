package timer

import "time"

type Timer struct {
	C <-chan time.Duration

	duration time.Duration
	c        chan time.Duration
	ticker   *time.Ticker
}

func New() *Timer {
	c := make(chan time.Duration)

	t := &Timer{
		C: c,
		c: c,
	}

	go t.Run()

	return t
}

func (t *Timer) Run() {
	// wait for user to stop its environment before start the timer
	t.c <- t.duration

	t.ticker = time.NewTicker(time.Second)

	for {
		<-t.ticker.C

		t.duration += time.Second

		select {
		case t.c <- t.duration:
		default:
		}
	}
}

func (t *Timer) Stop() {
	t.ticker.Stop()
}

func (t *Timer) Start() {
	t.ticker.Reset(time.Second)
}
