package ui

import (
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

func podsTab() tab {
	return &listTab{
		mu:         &sync.Mutex{},
		headerTmps: podHeaders,
		itemType:   reflect.TypeOf(podItem{}),
	}
}

var podHeaders = []headerTmp{
	{"name", 3},
	{"status", 1},
	{"host", 2},
	{"container", 3},
	{"cont status", 1},
	{"started at", 1},
	{"restarts", 1},
}

type podItem struct {
	p api.Pod
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
	lhost := label(pr.p.Spec.Host)

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
			lcontStarted.Text = c.State.Termination.StartedAt.String()
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

func (p *podItem) setData(d interface{})      { p.p = *d.(*api.Pod) }
func (p podItem) sameData(d interface{}) bool { return p.p.Name == (*d.(*api.Pod)).Name }
func (p podItem) less(i listItem) bool        { return p.p.Name < i.(*podItem).p.Name }
