package ui

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/gizak/termui"
)

type rcsTab struct {
	log *log.Logger

	mu  *sync.Mutex
	rcs []rc
}

func (rt *rcsTab) Len() int           { return len(rt.rcs) }
func (rt *rcsTab) Less(i, j int) bool { return rt.rcs[i].rc.Name < rt.rcs[j].rc.Name }
func (rt *rcsTab) Swap(i, j int)      { rt.rcs[i], rt.rcs[j] = rt.rcs[j], rt.rcs[i] }

func (rt *rcsTab) update(e Event) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	p := *e.Data.(*api.ReplicationController)
	switch e.Type {
	case watch.Added:
		rt.rcs = append(rt.rcs, rc{rc: p})
	case watch.Modified:
		found := false
		for i, up := range rt.rcs {
			if up.rc.Name == p.Name {
				found = true
				rt.rcs[i].rc = p
				break
			}
		}
		if !found {
			rt.rcs = append(rt.rcs, rc{rc: p})
		}
	case watch.Deleted:
		for i, up := range rt.rcs {
			if up.rc.Name == p.Name {
				rt.Swap(i, rt.Len()-1)
				rt.rcs = rt.rcs[:rt.Len()-1]
				break
			}
		}
	}
	sort.Sort(rt)
}

func (rt *rcsTab) toRows() []*termui.Row {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rows := make([]*termui.Row, 0)

	// header
	lname := header("name")
	lreplicas := header("replicas")
	ltmpl := header("template")
	lselect := header("selectors")

	rows = append(rows, termui.NewRow(
		termui.NewCol(4, 0, lname),
		termui.NewCol(1, 0, lreplicas),
		termui.NewCol(3, 0, ltmpl),
		termui.NewCol(4, 0, lselect),
	))
	for _, r := range rt.rcs {
		rows = append(rows, r.toRows()...)
	}
	return rows
}

type rc struct {
	rc api.ReplicationController
}

func (r rc) toRows() []*termui.Row {
	lname := label(r.rc.Name)
	lreplicas := label(fmt.Sprintf("%d/%d", r.rc.Status.Replicas, r.rc.Spec.Replicas))
	if r.rc.Status.Replicas == r.rc.Spec.Replicas {
		lreplicas.TextFgColor = termui.ColorGreen
	} else {
		lreplicas.TextFgColor = termui.ColorYellow
	}
	var tmplName string
	if r.rc.Spec.Template != nil {
		tmplName = r.rc.Spec.Template.Name
	}
	if r.rc.Spec.TemplateRef != nil {
		tmplName = r.rc.Spec.TemplateRef.Name
	}
	ltmpl := label(tmplName)

	names := make([]string, 0, len(r.rc.Spec.Selector))
	for n := range r.rc.Spec.Selector {
		names = append(names, n)
	}
	sort.Strings(names)

	rows := make([]*termui.Row, 0, len(names))
	for i, n := range names {
		if i > 0 {
			lname = label("")
			lreplicas = label("")
			ltmpl = label("")
		}

		lselect := label(fmt.Sprintf("%s=%s", n, r.rc.Spec.Selector[n]))
		rows = append(rows, termui.NewRow(
			termui.NewCol(4, 0, lname),
			termui.NewCol(1, 0, lreplicas),
			termui.NewCol(3, 0, ltmpl),
			termui.NewCol(4, 0, lselect),
		))
	}

	return rows
}
