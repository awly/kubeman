package ui

import (
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

func podsTab(ui *UI) tab {
	return &listTab{
		ui:       ui,
		mu:       &sync.Mutex{},
		headers:  podHeaders,
		itemType: reflect.TypeOf(podItem{}),
	}
}

var podHeaders = []header{
	{"name", 3},
	{"status", 1},
	{"host", 2},
	{"container", 3},
	{"cont status", 1},
	{"started at", 1},
	{"restarts", 1},
}

type podItem struct {
	p  api.Pod
	ui *UI
	lt *logTab
}

func (pr podItem) toRows() []*termui.Row {
	lname := label(pr.p.Name)
	lstatus := label(string(pr.p.Status.Phase))
	switch pr.p.Status.Phase {
	case api.PodRunning, api.PodSucceeded:
		lstatus.TextFgColor = termui.ColorGreen
	case api.PodPending:
		lstatus.TextFgColor = termui.ColorYellow
	case api.PodFailed, api.PodUnknown:
		lstatus.TextFgColor = termui.ColorRed
	}
	lhost := label(pr.p.Status.HostIP)
	if pr.p.Status.HostIP == "" {
		lhost.Text = pr.p.Status.PodIP
	}

	rows := make([]*termui.Row, 0, len(pr.p.Spec.Containers))
	for i, c := range pr.p.Status.ContainerStatuses {
		if i > 0 {
			lname = label("")
			lstatus = label("")
			lhost = label("")
		}
		lcont := label(c.Image)
		lcontStatus := label("")
		lcontStarted := label("")
		switch {
		case c.State.Running != nil:
			lcontStatus.Text = "Running"
			lcontStatus.TextFgColor = termui.ColorGreen
			lcontStarted.Text = c.State.Running.StartedAt.Format(time.Stamp)
		case c.State.Terminated != nil:
			lcontStatus.Text = "Terminated"
			lcontStatus.TextFgColor = termui.ColorRed
			lcontStarted.Text = c.State.Terminated.StartedAt.String()
		default:
			lcontStatus.Text = "Waiting"
			lcontStatus.TextFgColor = termui.ColorYellow
		}
		lrestarts := label(strconv.Itoa(c.RestartCount))

		rows = append(rows, termui.NewRow(
			termui.NewCol(3, 0, lname),
			termui.NewCol(1, 0, lstatus),
			termui.NewCol(2, 0, lhost),
			termui.NewCol(3, 0, lcont),
			termui.NewCol(1, 0, lcontStatus),
			termui.NewCol(1, 0, lcontStarted),
			termui.NewCol(1, 0, lrestarts),
		))
	}
	return rows
}

func (p *podItem) init(ui *UI)                { p.ui = ui }
func (p *podItem) setData(d interface{})      { p.p = *d.(*api.Pod) }
func (p podItem) sameData(d interface{}) bool { return p.p.Name == (*d.(*api.Pod)).Name }
func (p podItem) less(i listItem) bool        { return p.p.Name < i.(*podItem).p.Name }

func (p *podItem) handleEvent(e termui.Event) {
	switch e.Type {
	case termui.EventKey:
		switch e.Ch {
		case 'S':
			if err := p.ui.api.StopPod(p.p.Name); err != nil {
				log.Println(err)
			}
		case 'l':
			if p.lt != nil {
				p.ui.selected = "pods"
				p.ui.body = p.ui.tabs["pods"]
				p.lt.clean()
			} else {
				p.streamLogs()
			}
		}
	}
}

func (pi *podItem) streamLogs() {
	s, err := pi.ui.api.Logs(pi.p.Name, pi.p.Spec.Containers[0].Name, true)
	if err != nil {
		log.Println(err)
		return
	}
	lt := &logTab{
		in:        s,
		height:    termui.TermHeight() - 2,
		redraw:    pi.ui.redrawBody,
		uiUpdatef: pi.handleEvent,
	}
	lt.cleanf = func() {
		delete(pi.ui.tabs, "logs")
		pi.lt = nil
	}
	pi.lt = lt
	pi.ui.body = lt
	go lt.stream()
}
