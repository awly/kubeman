package ui

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/gizak/termui"
)

type Event struct {
	Resource string
	Type     watch.EventType
	Data     interface{}
}

func (ui *UI) eventLoop(ec <-chan termui.Event) {
	for e := range ec {
		ui.log.Printf("%+v", e)
		switch e.Type {
		case termui.EventInterrupt:
			close(ui.exitch)
			return
		case termui.EventResize:
			ui.mu.Lock()
			termui.Body.Width = termui.TermWidth()
			ui.mu.Unlock()
			ui.RedrawTabs()
			ui.RedrawBody()
		case termui.EventError:
			close(ui.exitch)
			return
		case termui.EventKey:
			if e.Key == termui.KeyCtrlC || e.Ch == 'q' {
				close(ui.exitch)
				return
			}
			if e.Ch >= '1' && e.Ch <= '9' {
				i := e.Ch - '1'
				tabs := ui.tabNames()
				if int(i) >= len(tabs) {
					continue
				}
				ui.SelectTab(ui.tabNames()[i])
				continue
			}
		case termui.EventMouse:
			// Top 2 rows are tabs
			if e.MouseY < 2 {
				// Tab index = X / tabWidth
				i := e.MouseX / (termui.TermWidth() / len(ui.tabs))
				ui.SelectTab(ui.tabNames()[i])
			}
		}
	}
}

func (ui *UI) SelectTab(name string) {
	ui.mu.Lock()
	ui.selected = name
	ui.mu.Unlock()
	ui.RedrawTabs()
	ui.RedrawBody()
}

func (ui *UI) updateLoop() {
	for e := range ui.Updates {
		handleUpdate(ui, e)
	}
}

func handleUpdate(ui *UI, e Event) {
	ui.log.Printf("%+v", e)
	if e.Type == watch.Error {
		return
	}
	t, ok := ui.tabs[e.Resource]
	if !ok {
		ui.log.Println("unsupported resource type", e.Resource)
		return
	}
	t.update(e)
	ui.RedrawBody()
}
