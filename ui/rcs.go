package ui

import (
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/alytvynov/kubeman/client"
	"github.com/gizak/termui"
)

func rcsTab() tab {
	return &listTab{
		mu:         &sync.Mutex{},
		headerTmps: rcHeaders,
		itemType:   reflect.TypeOf(rcItem{}),
	}
}

var rcHeaders = []headerTmp{
	{"name", 4},
	{"replicas", 1},
	{"template", 3},
	{"selectors", 4},
}

type rcItem struct {
	rc api.ReplicationController
}

func (r rcItem) toRows() []*termui.Row {
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

func (p *rcItem) setData(d interface{}) { p.rc = *d.(*api.ReplicationController) }
func (p rcItem) sameData(d interface{}) bool {
	return p.rc.Name == (*d.(*api.ReplicationController)).Name
}
func (p rcItem) less(i listItem) bool { return p.rc.Name < i.(*rcItem).rc.Name }

func (p rcItem) handleEvent(c *client.Client, e termui.Event) {}
