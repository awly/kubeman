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

type servicesTab struct {
	log *log.Logger

	mu       *sync.Mutex
	services []service
}

func (st *servicesTab) Len() int           { return len(st.services) }
func (st *servicesTab) Less(i, j int) bool { return st.services[i].s.Name < st.services[j].s.Name }
func (st *servicesTab) Swap(i, j int) {
	st.services[i], st.services[j] = st.services[j], st.services[i]
}

func (st *servicesTab) update(e Event) {
	st.mu.Lock()
	defer st.mu.Unlock()
	p := *e.Data.(*api.Service)
	switch e.Type {
	case watch.Added:
		st.services = append(st.services, service{s: p})
	case watch.Modified:
		found := false
		for i, up := range st.services {
			if up.s.Name == p.Name {
				found = true
				st.services[i].s = p
				break
			}
		}
		if !found {
			st.services = append(st.services, service{s: p})
		}
	case watch.Deleted:
		for i, up := range st.services {
			if up.s.Name == p.Name {
				st.Swap(i, st.Len()-1)
				st.services = st.services[:st.Len()-1]
				break
			}
		}
	}
	sort.Sort(st)
}

func (st *servicesTab) toRows() []*termui.Row {
	st.mu.Lock()
	defer st.mu.Unlock()
	rows := make([]*termui.Row, 0, len(st.services)+1)

	// header
	lname := label("name")
	lname.TextFgColor = termui.ColorWhite | termui.AttrBold
	ltype := label("type")
	ltype.TextFgColor = termui.ColorWhite | termui.AttrBold
	lip := label("ip")
	lip.TextFgColor = termui.ColorWhite | termui.AttrBold
	lport := label("port")
	lport.TextFgColor = termui.ColorWhite | termui.AttrBold

	rows = append(rows, termui.NewRow(
		termui.NewCol(3, 0, lname),
		termui.NewCol(1, 0, ltype),
		termui.NewCol(1, 0, lip),
		termui.NewCol(1, 0, lport),
	))
	for _, p := range st.services {
		rows = append(rows, p.toRows()...)
	}
	return rows
}

type service struct {
	s api.Service
}

func (s service) toRows() []*termui.Row {
	lname := label(s.s.Name)
	ltype := label(string(s.s.Spec.Type))
	lip := label(s.s.Spec.PortalIP)

	// TODO: seems like LoadBalancer field is WIP right now and is not
	// populated. Remove the len check later.
	if s.s.Spec.Type == api.ServiceTypeLoadBalancer && len(s.s.Status.LoadBalancer.Ingress) > 0 {
		lip.Text = s.s.Status.LoadBalancer.Ingress[0].IP
	}

	rows := make([]*termui.Row, 0, len(s.s.Spec.Ports))
	for i, p := range s.s.Spec.Ports {
		if i > 0 {
			lname = label("")
			ltype = label("")
			lip = label("")
		}
		lport := label(fmt.Sprintf("%d -> %s", p.Port, p.TargetPort.String()))
		rows = append(rows, termui.NewRow(
			termui.NewCol(3, 0, lname),
			termui.NewCol(1, 0, ltype),
			termui.NewCol(1, 0, lip),
			termui.NewCol(1, 0, lport),
		))
	}
	return rows
}
