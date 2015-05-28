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

	go ui.renderLoop()
	go ui.updateLoop()
	go ui.eventLoop()

	return ui, nil
}

func (ui *UI) renderLoop() {
	for range time.NewTicker(time.Second).C {
		termui.Render(termui.Body)
	}
}

func (ui *UI) updateLoop() {
	for e := range ui.Updates {
		ui.log.Println(e)
	}
}

func (ui *UI) eventLoop() {
	for e := range termui.EventCh() {
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
