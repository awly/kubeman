package ui

import (
	"log"
	"sort"

	"github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

type UI struct {
	Updates chan Event

	tabs     map[string]tab
	selected string
	log      *log.Logger
	exitch   chan struct{}
}

type tab interface {
	update(Event)
	toRows() []*termui.Row
}

func New(l *log.Logger) (*UI, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}
	termbox.SetInputMode(termbox.InputCurrent | termbox.InputMouse)

	uc := make(chan Event)
	exitch := make(chan struct{})
	ui := &UI{
		Updates: uc,
		log:     l,
		exitch:  exitch,
		tabs: map[string]tab{
			"pods":     &podsTab{log: l},
			"services": &servicesTab{log: l},
		},
		selected: "pods",
	}

	ui.RedrawTabs()
	ui.RedrawBody()

	go ui.updateLoop()
	go ui.eventLoop()

	return ui, nil
}

func (ui *UI) RedrawTabs() {
	// Tabs
	ui.log.Println(ui.selected)
	names := ui.tabNames()
	tabCols := make([]*termui.Row, 0, len(names))
	for _, n := range names {
		l := label(" " + n + " ")
		l.BgColor = termui.ColorBlue
		if n == ui.selected {
			l.BgColor = termui.ColorCyan
			l.TextBgColor = termui.ColorCyan
			l.TextFgColor = termui.ColorDefault | termui.AttrBold
		}
		l.Height = 2
		l.PaddingLeft = 1
		tabCols = append(tabCols, termui.NewCol(12/len(ui.tabs), 0, l))
	}
	if len(termui.Body.Rows) > 1 {
		termui.Body.Rows[0] = termui.NewRow(tabCols...)
	} else {
		termui.Body.AddRows(termui.NewRow(tabCols...))
	}
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (ui *UI) updateLoop() {
	for e := range ui.Updates {
		handleUpdate(ui, e)
	}
}

func (ui *UI) eventLoop() {
	ec := termui.EventCh()
	for e := range ec {
		ui.log.Printf("%+v", e)
		switch e.Type {
		case termui.EventInterrupt:
			close(ui.exitch)
		case termui.EventResize:
			termui.Body.Width = termui.TermWidth()
			ui.RedrawTabs()
			ui.RedrawBody()
		case termui.EventError:
			close(ui.exitch)
		case termui.EventKey:
			if e.Key == termui.KeyCtrlC || e.Ch == 'q' {
				close(ui.exitch)
			}
		case termui.EventMouse:
			// Top 2 rows are tabs
			if e.MouseY < 2 {
				// Tab index = X / tabWidth
				i := e.MouseX / (termui.TermWidth() / len(ui.tabs))

				ui.selected = ui.tabNames()[i]
				ui.RedrawTabs()
				ui.RedrawBody()
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

func (ui *UI) RedrawBody() {
	termui.Body.Rows = append(termui.Body.Rows[:1], ui.tabs[ui.selected].toRows()...)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func label(text string) *termui.Par {
	l := termui.NewPar(text)
	l.Height = 1
	l.HasBorder = false
	return l
}

func (ui *UI) tabNames() []string {
	names := make([]string, 0, len(ui.tabs))
	for n := range ui.tabs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
