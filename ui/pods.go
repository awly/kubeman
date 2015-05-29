package ui

import (
	"log"
	"sort"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/gizak/termui"
)

type podsTab struct {
	log  *log.Logger
	pods []pod
}

func (pl *podsTab) Len() int           { return len(pl.pods) }
func (pl *podsTab) Less(i, j int) bool { return pl.pods[i].p.Name < pl.pods[j].p.Name }
func (pl *podsTab) Swap(i, j int)      { pl.pods[i], pl.pods[j] = pl.pods[j], pl.pods[i] }

func (pl *podsTab) update(e Event) {
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
	rows := make([]*termui.Row, 0)

	// header
	lname := header("name")
	lstatus := header("status")
	lhost := header("host")
	lcont := header("container")

	rows = append(rows, termui.NewRow(
		termui.NewCol(3, 0, lname),
		termui.NewCol(1, 0, lstatus),
		termui.NewCol(2, 0, lhost),
		termui.NewCol(2, 0, lcont),
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
	lhost := label(pr.p.Spec.Host)

	rows := make([]*termui.Row, 0, len(pr.p.Spec.Containers))
	for i, c := range pr.p.Spec.Containers {
		if i != 0 {
			lname.Text = ""
			lstatus.Text = ""
			lhost.Text = ""
		}
		lcont := label(c.Image)
		rows = append(rows, termui.NewRow(
			termui.NewCol(3, 0, lname),
			termui.NewCol(1, 0, lstatus),
			termui.NewCol(2, 0, lhost),
			termui.NewCol(2, 0, lcont),
		))
	}
	return rows
}
