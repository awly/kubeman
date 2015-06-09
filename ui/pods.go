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
	{"container", 2},
	{"cont status", 1},
	{"started at", 2},
	{"restarts", 1},
}

type podItem struct {
	p  api.Pod
	ui *UI
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
	lhost := label(pr.p.Spec.NodeName)

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
		case c.State.Termination != nil:
			lcontStatus.Text = "Terminated"
			lcontStatus.TextFgColor = termui.ColorRed
			lcontStarted.Text = c.State.Termination.StartedAt.Format(time.Stamp)
		default:
			lcontStatus.Text = "Waiting"
			lcontStatus.TextFgColor = termui.ColorYellow
		}
		lrestarts := label(strconv.Itoa(c.RestartCount))

		rows = append(rows, termui.NewRow(
			termui.NewCol(3, 0, lname),
			termui.NewCol(1, 0, lstatus),
			termui.NewCol(2, 0, lhost),
			termui.NewCol(2, 0, lcont),
			termui.NewCol(1, 0, lcontStatus),
			termui.NewCol(2, 0, lcontStarted),
			termui.NewCol(1, 0, lrestarts),
		))
	}
	if len(rows) == 0 {
		rows = append(rows, termui.NewRow(
			termui.NewCol(3, 0, lname),
			termui.NewCol(1, 0, lstatus),
			termui.NewCol(2, 0, lhost),
			termui.NewCol(2, 0, label("")),
			termui.NewCol(1, 0, label("")),
			termui.NewCol(2, 0, label("")),
			termui.NewCol(1, 0, label("")),
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
			showLogTab(p.ui, p.p)
		}
	}
}
