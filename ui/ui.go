package ui

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

type UI struct {
	Updates chan Event

	tabs     map[string]tab
	selected string
	log      *log.Logger
	exitch   chan struct{}

	// protects selected and termui.Body
	mu *sync.Mutex
}

func New(l *log.Logger) (*UI, error) {
	uc := make(chan Event)
	exitch := make(chan struct{})
	ui := &UI{
		Updates: uc,
		log:     l,
		exitch:  exitch,
		tabs: map[string]tab{
			"pods":     podsTab(),
			"services": servicesTab(),
			"rcs":      rcsTab(),
			"nodes":    nodesTab(),
		},
		selected: "pods",
		mu:       &sync.Mutex{},
	}

	go ui.updateLoop()
	go ui.eventLoop(termui.EventCh())

	if err := termui.Init(); err != nil {
		return nil, err
	}
	termbox.SetInputMode(termbox.InputCurrent | termbox.InputMouse)

	ui.RedrawTabs()
	ui.RedrawBody()

	return ui, nil
}

func (ui *UI) RedrawTabs() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	names := ui.tabNames()
	tabCols := make([]*termui.Row, 0, len(names))
	spaceCols := make([]*termui.Row, 0, len(names))
	for i, n := range names {
		l := label(fmt.Sprintf(" %d: %s ", i+1, n))
		l.TextBgColor = termui.ColorBlue
		s := label("")
		s.PaddingRight = 1
		s.BgColor = termui.ColorBlue
		if n == ui.selected {
			l.Text += "* "
			l.TextFgColor = termui.ColorDefault | termui.AttrBold
			l.TextBgColor = termui.ColorCyan
			s.BgColor = termui.ColorCyan
		}
		tabCols = append(tabCols, termui.NewCol(12/len(ui.tabs), 0, l))
		spaceCols = append(spaceCols, termui.NewCol(12/len(ui.tabs), 0, s))
	}
	if len(termui.Body.Rows) > 2 {
		termui.Body.Rows[0] = termui.NewRow(tabCols...)
		termui.Body.Rows[1] = termui.NewRow(spaceCols...)
	} else {
		termui.Body.AddRows(termui.NewRow(tabCols...), termui.NewRow(spaceCols...))
	}
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (ui *UI) RedrawBody() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	termui.Body.Rows = append(termui.Body.Rows[:2], ui.tabs[ui.selected].toRows()...)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (ui *UI) Close() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	termui.Close()
}

func (ui *UI) ExitCh() <-chan struct{} {
	return ui.exitch
}

func (ui *UI) tabNames() []string {
	names := make([]string, 0, len(ui.tabs))
	for n := range ui.tabs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func label(text string) *termui.Par {
	l := termui.NewPar(text)
	l.Height = 1
	l.HasBorder = false
	return l
}

func header(text string) *termui.Par {
	l := label(text)
	l.TextFgColor = termui.ColorWhite | termui.AttrBold
	l.PaddingLeft = 1
	return l
}
