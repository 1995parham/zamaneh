package app

import (
	"context"
	"time"

	"github.com/1995parham/zamaneh/timer"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/sirupsen/logrus"
)

const (
	defaultFg           = cell.ColorDefault
	RefreshRate         = 100 * time.Millisecond
	PaddingLeftPercent  = 5
	PaddingRightPercent = 5
)

type App struct {
	Term    terminalapi.Terminal
	Context context.Context
	Cancel  func()

	isStop bool
	timer  *timer.Timer

	update chan []cell.Option
	cron   *segmentdisplay.SegmentDisplay
	panel  *container.Container
	topic  string
}

func New(topic string) (*App, error) {
	t, err := termbox.New()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		Term:    t,
		Context: ctx,
		Cancel:  cancel,

		timer: timer.New(),

		update: make(chan []cell.Option),
		topic:  topic,
	}, nil
}

func (a *App) Run() error {
	if err := a.build(); err != nil {
		return err
	}

	if err := a.layout(); err != nil {
		return err
	}

	go a.cronner()

	if err := termdash.Run(
		a.Context, a.Term, a.panel,
		termdash.KeyboardSubscriber(a.keyHandler),
		termdash.RedrawInterval(RefreshRate),
	); err != nil {
		return err
	}

	return nil
}

func (a *App) build() error {
	cron, err := segmentdisplay.New()
	if err != nil {
		return err
	}

	a.cron = cron

	if err := cron.Write([]*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(
			"00:00:00",
		),
	}); err != nil {
		return err
	}

	return nil
}

func (a *App) layout() error {
	c, err := container.New(
		a.Term,
		container.Border(linestyle.Light),
		container.BorderTitle(a.topic),
		container.PlaceWidget(a.cron),
		container.PaddingLeftPercent(PaddingLeftPercent),
		container.PaddingRightPercent(PaddingRightPercent),
	)
	if err != nil {
		return err
	}

	a.panel = c

	return nil
}

func (a *App) keyHandler(k *terminalapi.Keyboard) {
	if k.Key == 'q' || k.Key == 'Q' {
		a.Close()
	}

	if k.Key == keyboard.KeySpace {
		if a.isStop {
			a.timer.Start()
			a.update <- []cell.Option{}
			a.isStop = false
		} else {
			a.timer.Stop()
			a.update <- []cell.Option{
				cell.FgColor(cell.ColorCyan),
			}
			a.isStop = true
		}
	}
}

func (a *App) cronner() {
	var d time.Duration

	defaultOptions := []cell.Option{
		cell.FgColor(defaultFg),
	}
	opts := defaultOptions

	for {
		select {
		case d = <-a.timer.C:
		case opts = <-a.update:
			if len(opts) == 0 {
				opts = defaultOptions
			}
		}

		if err := a.cron.Write([]*segmentdisplay.TextChunk{
			segmentdisplay.NewChunk(
				time.Time{}.Add(d).Format("15:04:05"),
				segmentdisplay.WriteCellOpts(opts...),
			),
		}); err != nil {
			logrus.Fatal(err)
		}
	}
}

func (a *App) Close() {
	a.Cancel()
	a.Term.Close()
}
