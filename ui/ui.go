package ui

import (
	"log"
	"time"

	"github.com/gizak/termui"
)

type UI struct {
	Updates chan Event
	log     *log.Logger
	exitch  chan struct{}
}

func New(l *log.Logger) (*UI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	uc := make(chan Event)
	exitch := make(chan struct{})
	ui := &UI{
		Updates: uc,
		log:     l,
		exitch:  exitch,
	}

	buildLayout()

	go ui.renderLoop()
	go ui.updateLoop()
	go ui.eventLoop()

	return ui, nil
}

func buildLayout() {
	// Tabs
	lpods := label("pods")
	lpods.Height = 2
	lrcs := label("rcs")
	lrcs.Height = 2
	lservices := label("services")
	lservices.Height = 2
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(2, 2, lpods),
			termui.NewCol(2, 2, lrcs),
			termui.NewCol(2, 2, lservices),
		),
	)

	// Content

	termui.Body.Align()
}

func (ui *UI) renderLoop() {
	for range time.NewTicker(100 * time.Millisecond).C {
		termui.Render(termui.Body)
	}
}

func (ui *UI) updateLoop() {
	for e := range ui.Updates {
		ui.log.Println(e)
	}
}

func (ui *UI) eventLoop() {
	ec := termui.EventCh()
	for e := range ec {
		ui.log.Printf("%+v", e)
		switch e.Type {
		case termui.EventInterrupt:
			close(ui.exitch)
		case termui.EventError:
			close(ui.exitch)
		case termui.EventKey:
			if e.Key == termui.KeyCtrlC || e.Ch == 'q' {
				close(ui.exitch)
			}
		}
	}
}

func (ui *UI) Close() {
	termui.Close()
}

func (ui *UI) ExitCh() <-chan struct{} {
	return ui.exitch
}

type Event struct {
	Resource string
	Type     string
	Data     interface{}
}

func label(text string) *termui.Par {
	l := termui.NewPar(text)
	l.Height = 1
	l.HasBorder = false
	return l
}
