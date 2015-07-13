package ui

import (
	"fmt"
	"sort"
	"sync"

	"github.com/alytvynov/kubeman/client"
	"github.com/alytvynov/termui"
)

type UI struct {
	tabs     map[string]tab
	selected string
	body     tab
	status   tab
	exitch   chan struct{}
	api      *client.Client

	// protects selected and termui.Body
	mu *sync.Mutex
}

func New(c *client.Client) (*UI, error) {
	exitch := make(chan struct{})
	ui := &UI{
		exitch:   exitch,
		selected: "pods",
		mu:       &sync.Mutex{},
		api:      c,
	}
	ui.status = &statusBar{ui: ui}
	ui.tabs = map[string]tab{
		"pods":     podsTab(ui),
		"services": servicesTab(ui),
		"rcs":      rcsTab(ui),
		"nodes":    nodesTab(ui),
	}
	ui.body = ui.tabs[ui.selected]

	go ui.eventLoop(termui.EventCh())

	if err := termui.Init(); err != nil {
		return nil, err
	}

	termui.Body.Rows = make([]*termui.Row, termui.TermHeight())
	for i := range termui.Body.Rows {
		termui.Body.Rows[i] = termui.NewRow(termui.NewCol(12, 0, label("")))
	}

	ui.redrawTabs()
	ui.redrawBody()

	go ui.watchUpdates()

	return ui, nil
}

func (ui *UI) redrawTabs() {
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
	termui.Body.Rows[0] = termui.NewRow(tabCols...)
	termui.Body.Rows[1] = termui.NewRow(spaceCols...)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (ui *UI) redrawBody() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	n := copy(termui.Body.Rows[2:len(termui.Body.Rows)-1], ui.body.toRows())
	for i := n + 2; i < len(termui.Body.Rows)-1; i++ {
		termui.Body.Rows[i] = termui.NewRow(termui.NewCol(12, 0, label("")))
	}
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (ui *UI) redrawStatus() {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	termui.Body.Rows[len(termui.Body.Rows)-1] = ui.status.toRows()[0]
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
