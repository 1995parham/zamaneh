package app

import (
	"context"
	"fmt"
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
	"github.com/mum4k/termdash/widgets/text"
	"github.com/sirupsen/logrus"
)

const (
	yellowFg = 226
	orangeFg = 215
	cyanFg   = 123
	pinkFg   = 218

	RefreshRate    = 100 * time.Millisecond
	PaddingPercent = 5
)

type App struct {
	Term    terminalapi.Terminal
	Context context.Context
	Cancel  func()

	isStop bool
	timer  *timer.Timer

	update chan []cell.Option
	cron   *segmentdisplay.SegmentDisplay
	text   *text.Text
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
			segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorNumber(yellowFg))),
		),
	}); err != nil {
		return err
	}

	txt, err := text.New()
	if err != nil {
		return err
	}

	a.text = txt

	if err := txt.Write(
		fmt.Sprintf(
			"Since: %s\nTopic: %s",
			time.Now().Format(time.RFC822),
			a.topic,
		),
		text.WriteCellOpts(
			cell.FgColor(cell.ColorNumber(orangeFg)),
		),
	); err != nil {
		return err
	}

	if err := txt.Write(
		"\n\nSpending time with you is so precious,\nI love every minute that we are together.",
		text.WriteCellOpts(
			cell.FgColor(cell.ColorNumber(pinkFg)),
		),
	); err != nil {
		return err
	}

	return nil
}

func (a *App) layout() error {
	c, err := container.New(
		a.Term,
		container.Border(linestyle.Light),
		container.SplitHorizontal(
			container.Top(
				container.PlaceWidget(a.cron),
			),
			container.Bottom(
				container.SplitVertical(
					container.Left(
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorCyan),
						container.BorderTitle("notes"),
						container.PlaceWidget(a.text),
						container.PaddingLeftPercent(PaddingPercent),
						container.PaddingRightPercent(PaddingPercent),
						container.PaddingTopPercent(PaddingPercent),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorGreen),
						container.BorderTitle("announcements"),
						container.PaddingLeftPercent(PaddingPercent),
						container.PaddingRightPercent(PaddingPercent),
						container.PaddingTopPercent(PaddingPercent),
					),
				),
			),
		),
		container.PaddingLeftPercent(PaddingPercent),
		container.PaddingRightPercent(PaddingPercent),
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
				cell.FgColor(cell.ColorNumber(cyanFg)),
			}
			a.isStop = true
		}
	}
}

func (a *App) cronner() {
	var d time.Duration

	defaultOptions := []cell.Option{
		cell.FgColor(cell.ColorNumber(yellowFg)),
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
			logrus.Error(err)
		}
	}
}

func (a *App) Close() {
	a.Cancel()
	a.Term.Close()
}
