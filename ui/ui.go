package ui

import (
	"log"
	"sort"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

type UI struct {
	Updates chan Event
	pods    podList
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
	ui.Redraw()

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

func (ui *UI) updateLoop() {
	for e := range ui.Updates {
		handleUpdate(ui, e)
	}
}

func (ui *UI) SetPods(pods []api.Pod) {
	upods := podList{pods: make([]pod, 0, len(pods))}
	for _, p := range pods {
		upods.pods = append(upods.pods, pod{p: p})
	}
	sort.Sort(upods)
	ui.pods = upods
	ui.Redraw()
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

func (ui *UI) Redraw() {
	termui.Body.Rows = append(termui.Body.Rows[:1], ui.pods.toRows()...)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func label(text string) *termui.Par {
	l := termui.NewPar(text)
	l.Height = 1
	l.HasBorder = false
	return l
}
