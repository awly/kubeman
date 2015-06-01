package ui

import (
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/gizak/termui"
)

type podsTab struct {
	log *log.Logger

	mu   *sync.Mutex
	pods []pod
}

func (pl *podsTab) Len() int           { return len(pl.pods) }
func (pl *podsTab) Less(i, j int) bool { return pl.pods[i].p.Name < pl.pods[j].p.Name }
func (pl *podsTab) Swap(i, j int)      { pl.pods[i], pl.pods[j] = pl.pods[j], pl.pods[i] }

func (pl *podsTab) update(e Event) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	p := *e.Data.(*api.Pod)
	switch e.Type {
	case watch.Added:
		pl.pods = append(pl.pods, pod{p: p})
	case watch.Modified:
		found := false
		for i, up := range pl.pods {
			if up.p.Name == p.Name {
				found = true
				pl.pods[i].p = p
				break
			}
		}
		if !found {
			pl.pods = append(pl.pods, pod{p: p})
		}
	case watch.Deleted:
		for i, up := range pl.pods {
			if up.p.Name == p.Name {
				pl.Swap(i, pl.Len()-1)
				pl.pods = pl.pods[:pl.Len()-1]
				break
			}
		}
	}
	sort.Sort(pl)
}

func (pl *podsTab) toRows() []*termui.Row {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	rows := make([]*termui.Row, 0)

	// header
	lname := header("name")
	lstatus := header("status")
	lhost := header("host")
	lcont := header("container")
	lcontStatus := header("cont status")
	lcontStarted := header("started at")
	lrestarts := header("restarts")

	rows = append(rows, termui.NewRow(
		termui.NewCol(3, 0, lname),
		termui.NewCol(1, 0, lstatus),
		termui.NewCol(2, 0, lhost),
		termui.NewCol(3, 0, lcont),
		termui.NewCol(1, 0, lcontStatus),
		termui.NewCol(1, 0, lcontStarted),
		termui.NewCol(1, 0, lrestarts),
	))
	for _, p := range pl.pods {
		rows = append(rows, p.toRows()...)
	}
	return rows
}

type pod struct {
	p api.Pod
}

func (pr pod) toRows() []*termui.Row {
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
