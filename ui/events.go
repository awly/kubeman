package ui

import (
	"log"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/alytvynov/termui"
)

type Event struct {
	Resource string
	Type     watch.EventType
	Data     interface{}
}

func (ui *UI) eventLoop(ec <-chan termui.Event) {
	for e := range ec {
		log.Printf("%+v", e)
		switch e.Type {
		case termui.EventInterrupt:
			close(ui.exitch)
			return
		case termui.EventResize:
			ui.mu.Lock()
			termui.Body.Width = termui.TermWidth()
			ui.mu.Unlock()
			ui.redrawTabs()
			ui.redrawBody()
		case termui.EventError:
			close(ui.exitch)
			return
		case termui.EventKey:
			if e.Key == termui.KeyCtrlC || e.Ch == 'q' {
				close(ui.exitch)
				return
			}
			if e.Ch == 'R' {
				ui.api.DisconnectWatches()
				continue
			}
			if e.Ch >= '1' && e.Ch <= '9' {
				i := e.Ch - '1'
				tabs := ui.tabNames()
				if int(i) >= len(tabs) {
					continue
				}
				log.Println("select tab i:", i, "name:", tabs[i])
				ui.SelectTab(tabs[i])
				continue
			}
			ui.updateTabUI(e)
		}
	}
}

func (ui *UI) SelectTab(name string) {
	ui.mu.Lock()
	ui.body.clean()
	ui.selected = name
	ui.body = ui.tabs[ui.selected]
	ui.mu.Unlock()
	ui.redrawTabs()
	ui.redrawBody()
}

func (ui *UI) updateTabUI(e termui.Event) {
	ui.mu.Lock()
	ui.body.uiUpdate(e)
	ui.mu.Unlock()
}

func (ui *UI) handleUpdate(e Event) {
	log.Printf("%+v", e)
	if e.Type == watch.Error {
		return
	}
	if e.Resource == "status" {
		ui.status.dataUpdate(e)
		return
	}
	t, ok := ui.tabs[e.Resource]
	if !ok {
		log.Println("unsupported resource type", e.Resource)
		return
	}
	t.dataUpdate(e)
}

func (ui *UI) statusUpdate(msg string) {
	ui.handleUpdate(Event{Resource: "status", Data: statusUpdate{msg: msg}})
}
