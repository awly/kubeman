package ui

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

func servicesTab() tab {
	return &listTab{
		mu:         &sync.Mutex{},
		headerTmps: serviceHeaders,
		itemType:   reflect.TypeOf(serviceItem{}),
	}
}

var serviceHeaders = []headerTmp{
	{"name", 4},
	{"type", 2},
	{"ip", 3},
	{"port", 3},
}

type serviceItem struct {
	service api.Service
}

func (s serviceItem) toRows() []*termui.Row {
	lname := label(s.service.Name)
	ltype := label(string(s.service.Spec.Type))
	lip := label(s.service.Spec.PortalIP)

	// TODO: seems like LoadBalancer field is WIP right now and is not
	// populated. Remove the len check later.
	if s.service.Spec.Type == api.ServiceTypeLoadBalancer && len(s.service.Status.LoadBalancer.Ingress) > 0 {
		lip.Text = s.service.Status.LoadBalancer.Ingress[0].IP
	}

	rows := make([]*termui.Row, 0, len(s.service.Spec.Ports))
	for i, p := range s.service.Spec.Ports {
		if i > 0 {
			lname = label("")
			ltype = label("")
			lip = label("")
		}
		lport := label(fmt.Sprintf("%d -> %s", p.Port, p.TargetPort.String()))
		rows = append(rows, termui.NewRow(
			termui.NewCol(4, 0, lname),
			termui.NewCol(2, 0, ltype),
			termui.NewCol(3, 0, lip),
			termui.NewCol(3, 0, lport),
		))
	}
	return rows
}

func (p *serviceItem) setData(d interface{})      { p.service = *d.(*api.Service) }
func (p serviceItem) sameData(d interface{}) bool { return p.service.Name == (*d.(*api.Service)).Name }
func (p serviceItem) less(i listItem) bool        { return p.service.Name < i.(*serviceItem).service.Name }
